package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	log "github.com/gobuild/log"
	"github.com/howeyc/fsnotify"
)

const windows = "windows"

var verbose = flag.Bool("v", false, "show verbose")

func init() {
	log.SetFlags(0)
	if runtime.GOOS == windows {
		log.SetPrefix("fswatch >>> ")
	} else {
		log.SetPrefix("\033[32mfswatch\033[0m >>> ")
	}
}

/*
type pathWatch struct {
	Include   string `json:"include"`
	reInclude *regexp.Regexp
	Exclude   string `json:"exclude"`
	reExclude *regexp.Regexp
	Depth     int `json:"depth"`
}

type fsWatch struct {
	PathWatches []pathWatch
	Command     []string `json:"command"`
	Cmd         string   `json:"cmd"` // if empty will add prefix(bash -c) and replace Command
	Signal      string   `json:"signal"`
	KillAsGroup bool     `json:"killasgroup"`
}

func init() {
	fw := fsWatch{
		PathWatches: []pathWatch{
			pathWatch{Include: "./", Exclude: "\\.svn"},
		},
		Cmd: "ls -l",
	}
	data, err := yaml.Marshal(fw)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
*/

type gowatch struct {
	Paths     []string `json:"paths"`
	Depth     int      `json:"depth"`
	Exclude   []string `json:"exclude"`
	reExclude []*regexp.Regexp
	Include   []string `json:"include"`
	reInclude []*regexp.Regexp
	//bufdur    time.Duration `json:"-"`
	Command interface{} `json:"command"` // can be string or []string
	cmd     []string
	Env     map[string]string `json:"env"`

	AutoRestart     bool          `json:"autorestart"`
	RestartInterval time.Duration `json:"restart-interval"`
	KillSignal      string        `json:"kill-signal"`

	w       *fsnotify.Watcher
	modtime map[string]time.Time
	sig     chan string
	sigOS   chan os.Signal
}

// Check if file matches
func (t *gowatch) match(file string) bool {
	file = filepath.Base(file)
	for _, rule := range t.reExclude {
		if rule.MatchString(file) {
			return false
		}
	}
	for _, rule := range t.reInclude {
		if rule.MatchString(file) {
			return true
		}
	}
	return len(t.reInclude) == 0 // if empty include, then return true
}

// Add dir and children (recursively) to watcher
func (t *gowatch) watchDirAndChildren(path string, depth int) error {
	if err := t.w.Watch(path); err != nil {
		return err
	}
	baseNumSeps := strings.Count(path, string(os.PathSeparator))
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			base := info.Name()
			if base != "." && strings.HasPrefix(base, ".") { // ignore hidden dir
				return filepath.SkipDir
			}

			pathDepth := strings.Count(path, string(os.PathSeparator)) - baseNumSeps
			if pathDepth > depth {
				return filepath.SkipDir
			}
			if *verbose {
				fmt.Println(">>> watch dir: ", path)
			}
			if err := t.w.Watch(path); err != nil {
				return err
			}
		}
		return nil
	})
}

// Create a fsnotify fswatch
// Initial vars
func (t *gowatch) Watch() (err error) {
	if t.w, err = fsnotify.NewWatcher(); err != nil {
		return
	}
	for _, path := range t.Paths {
		// translate env-vars
		if err = t.watchDirAndChildren(os.ExpandEnv(path), t.Depth); err != nil {
			log.Fatal(err)
		}
	}
	t.modtime = make(map[string]time.Time)
	t.sig = make(chan string)
	for _, patten := range t.Exclude {
		t.reExclude = append(t.reExclude, regexp.MustCompile(patten))
	}
	for _, patten := range t.Include {
		t.reInclude = append(t.reInclude, regexp.MustCompile(patten))
	}

	t.sigOS = make(chan os.Signal, 1)
	signal.Notify(t.sigOS, syscall.SIGINT)

	go t.drainExec()
	t.drainEvent()
	return
}

// filter fsevent and send to t.sig
func (t *gowatch) drainEvent() {
	for {
		select {
		case err := <-t.w.Error:
			log.Warnf("watch error: %s", err)
		case <-t.sigOS:
			t.sig <- "EXIT"
		case eve := <-t.w.Event:
			log.Debug(eve)
			changed := t.IsfileChanged(eve.Name)
			if changed && t.match(eve.Name) {
				log.Info(eve)
				select {
				case t.sig <- "KILL":
				default:
				}
			}
		}
	}
}

// Use modified time to judge if file changed
func (t *gowatch) IsfileChanged(p string) bool {
	p = filepath.Clean(p)
	fi, err := os.Stat(p)
	if err != nil {
		return true // if file not exists, just return true
	}
	curr := fi.ModTime()
	defer func() { t.modtime[p] = curr }()
	modt, ok := t.modtime[p]
	return !ok || curr.After(modt.Add(time.Second))
}

func (t *gowatch) drainExec() {
	log.Println("command:", t.cmd)
	var msg string
	for {
		startTime := time.Now()
		cmd := t.cmd
		if len(cmd) == 0 {
			cmd = []string{"echo", "no command specified"}
		}
		cPrintf(cGREEN, "exec start")
		c := startCmd(cmd[0], cmd[1:]...)
		// Start to run command
		err := c.Start()
		if err != nil {
			cPrintf("35", err.Error())
		}
		// Wait until killed or finished
		select {
		case msg = <-t.sig:
			cPrintf(cYELLOW, "program terminated, signal(%s)", t.KillSignal)
			if err := killCmd(c, t.KillSignal); err != nil {
				log.Errorf("group kill: %v", err)
			}
			if msg == "EXIT" {
				os.Exit(1)
			}
			goto SKIP_WAITING
		case err = <-chanGo(c.Wait):
			if err != nil {
				cPrintf(cPURPLE, "program exited: %v", err)
			}
		}
		log.Infof("finish in %s", time.Since(startTime))

		// Whether to restart right now
		if t.AutoRestart {
			goto SKIP_WAITING
		}
		cPrintf("33", "-- wait signal --")
		if msg = <-t.sig; msg == "EXIT" {
			os.Exit(1)
		}
	SKIP_WAITING:
		if t.RestartInterval > 0 {
			log.Infof("restart after %s", t.RestartInterval)
		}
		time.Sleep(t.RestartInterval)
	}
}

const jsonconf = ".fswatch.json"

var (
	gw = &gowatch{
		Paths:           []string{"."},
		Depth:           2,
		Exclude:         []string{},
		Include:         []string{"\\.(go|py|php|java|cpp|h|rb)$"},
		AutoRestart:     false,
		RestartInterval: 0,
		KillSignal:      "KILL",
	}
	confExists = false
	extInclude string
)

// parse command flag
func flagParse() {
	gw.Env = map[string]string{"POWERD_BY": "github.com/codeskyblue/fswatch"}
	// load jsonconf
	if fd, err := os.Open(jsonconf); err == nil {
		if er := json.NewDecoder(fd).Decode(gw); er != nil {
			log.Fatalf("json decode error: %v", er)
		}
		for key, val := range gw.Env {
			os.Setenv(key, val)
		}
		confExists = true
	}
	flag.DurationVar(&gw.RestartInterval, "ri", gw.RestartInterval, "restart interval")
	flag.BoolVar(&gw.AutoRestart, "r", gw.AutoRestart, "auto restart")
	flag.StringVar(&gw.KillSignal, "k", gw.KillSignal, "kill signal")
	flag.StringVar(&extInclude, "ext", "", "extensions eg: [cpp,c,h]")
	flag.IntVar(&gw.Depth, "d", 2, "watch depth")
	flag.Parse()
}

func main() {
	flagParse()

	if len(os.Args) == 1 && !confExists {
		fmt.Printf("Create %s file [y/n]: ", strconv.Quote(jsonconf))
		var yn = "y"
		fmt.Scan(&yn)
		gw.Command = "echo helloworld"
		if strings.ToUpper(strings.TrimSpace(yn)) == "Y" {
			data, _ := json.MarshalIndent(gw, "", "    ")
			ioutil.WriteFile(jsonconf, data, 0644)
			fmt.Printf("use notepad++ or vim to edit %s\n", strconv.Quote(jsonconf))
		}
		return
	}

	if flag.NArg() > 0 {
		gw.Command = []string(flag.Args())
	}
	if extInclude != "" {
		for _, ext := range strings.Split(extInclude, ",") {
			gw.Include = append(gw.Include, "\\."+ext+"$")
		}
	}

	// []string is unmarshaled as []interface{}
	if v, ok := gw.Command.([]interface{}); ok {
		cmd := make([]string, 0, len(v))
		for _, i := range v {
			if s, ok := i.(string); ok {
				cmd = append(cmd, s)
			} else {
				log.Fatalf("check you config file. \"command\" must be string or []string, got %T", gw.Command)
			}
		}
		gw.Command = cmd
	}

	switch gw.Command.(type) {
	default:
		log.Fatalf("check you config file. \"command\" must be string or []string, got %T", gw.Command)
	case string:
		if runtime.GOOS == windows {
			gw.cmd = []string{"cmd", "/c", gw.Command.(string)}
		} else {
			gw.cmd = []string{"bash", "-c", gw.Command.(string)}
		}
	case []string:
		gw.cmd = gw.Command.([]string)
	}

	log.Fatal(gw.Watch())
}
