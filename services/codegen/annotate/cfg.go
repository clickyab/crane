package annotate

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	group = regexp.MustCompile(`(?ms)^\s*(/[*]|[/]{2}|)\s*@(\S+)\s*{([^}]+)}\s*$`)
	line  = regexp.MustCompile(`\s*([/]{2}\s*|)((\S+)\s*=(.*)|.*)`)
)

type Annotate struct {
	Items map[string]string
	Name  string
}

type AnnotateGroup []Annotate

func loadFromGroup(g string) (Annotate, error) {
	res := Annotate{Items: make(map[string]string)}
	lne := line.FindAllStringSubmatch(g, -1)
	for i := range lne {
		if len(lne[i]) == 5 {
			l := strings.Trim(lne[i][2], " /\n\t")
			k := strings.Trim(lne[i][3], " \n\t")
			v := strings.Trim(lne[i][4], " \n\t")

			if k != "" {
				res.Items[k] = v
			} else {
				if l != "" {
					return Annotate{}, fmt.Errorf("invalid line '%s'", l)
				}
			}
		}
	}

	return res, nil
}

func LoadFromComment(c string) (AnnotateGroup, error) {
	// First find groups
	var res AnnotateGroup
	grps := group.FindAllStringSubmatch(c, -1)

	for i := range grps {
		if len(grps[i]) == 4 {
			a, err := loadFromGroup(grps[i][3])
			if err != nil {
				return nil, err
			}
			a.Name = grps[i][2]
			res = append(res, a)
		}
	}

	return res, nil
}
