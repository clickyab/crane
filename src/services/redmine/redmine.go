package redmine

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"services/config"
	"services/safe"

	"github.com/Sirupsen/logrus"
	"github.com/fzerorubigd/go-redmine"
	"gopkg.in/fzerorubigd/onion.v2"
)

var (
	url            = config.RegisterString("services.redmine.url", "")
	apiKey         = config.RegisterString("services.redmine.apikey", "")
	newIssueTypeID = config.RegisterInt("services.redmine.new_issue_type_id", 0)
	projectID      = config.RegisterInt("services.redmine.project_id", 0)
	active         = config.RegisterBoolean("services.redmine.active", false)
)

type reporter struct {
}

func (reporter) Initialize(*onion.Onion) []onion.Layer {
	return nil
}

// Loaded is called after config loading, so the active is ready here
func (r *reporter) Loaded() {
	if *active {
		safe.Register(r)
	}
}

func (reporter) Recover(err error, ds []byte, extra ...interface{}) {
	c := redmine.NewClient(*url, *apiKey)

	// redmine can not accept more than 255 character title
	var title error
	if len(err.Error()) > 200 {
		str := err.Error()
		title = errors.New(str[:200] + "...")
	}

	stack := string(ds)
	for i := range extra {
		if t, ok := extra[i].(*http.Request); ok {
			if b, err := httputil.DumpRequest(t, true); err != nil {
				stack += "\n\n the https request dump : \n\n%s" + string(b)
				continue
			}
		}

		stack += fmt.Sprintf("Extra data :\n %T => %+v", extra[i], extra[i])
	}

	var filters []redmine.IssueFilter
	filters = append(filters, redmine.IssueFilter{Key: "limit", Value: "1"})
	filters = append(filters, redmine.IssueFilter{Key: "subject", Value: title.Error()})
	//filters = append(filters, redmine.IssueFilter{Key: "status_id", Value: "open"})

	issues, err := c.FilterIssues(filters...)
	if err != nil {
		logrus.Warn(err)
		return
	}
	var is *redmine.Issue
	if len(issues) > 0 {
		for i := range issues {
			if issues[i].Status.Id == *newIssueTypeID {
				is = &issues[i]
				break
			}
		}
	}

	if is != nil {
		is.Notes = stack
		err := c.UpdateIssue(*is)
		if err != nil {
			logrus.Warn(err)
		}
	} else {
		is = &redmine.Issue{}
		is.Subject = title.Error()
		is.Description = stack
		is.ProjectId = *projectID

		_, err := c.CreateIssue(*is)
		if err != nil {
			logrus.Warn(err)
		}
	}
}

func init() {
	config.Register(&reporter{})
}
