package main

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/url"
	"fmt"
)

type Build struct {
	BuildStatusPrevious   string   `json:"buildStatusPrevious,omitempty"`
	BuildId               string   `json:"buildId,omitempty"`
	TriggeredBy           string   `json:"triggeredBy,omitempty"`
	BuildStatus           string   `json:"buildStatus,omitempty"`
	AgentHostname         string   `json:"agentHostname,omitempty"`
	BuildFullName         string   `json:"buildFullName,omitempty"`
	BuildTypeId           string   `json:"buildTypeId,omitempty"`
	Message               string   `json:"message,omitempty"`
	Text                  string   `json:"text,omitempty"`
	NotifyType            string   `json:"notifyType,omitempty"`
	AgentName             string   `json:"agentName,omitempty"`
	BuildResult           string   `json:"buildResult,omitempty"`
	BuildRunner           string   `json:"buildRunner,omitempty"`
	ProjectId             string   `json:"projectId,omitempty"`
	ProjectName           string   `json:"projectName,omitempty"`
	AgentOs               string   `json:"agentOs,omitempty"`
	BuildNumber           string   `json:"buildNumber,omitempty"`
	BuildName             string   `json:"buildName,omitempty"`
}

var lastBuild = Build{}
var lastForm = url.Values{} //map[string]string{}

func main() {
	m := martini.New()

	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	// Setup routes
	r := martini.NewRouter()
	r.Get(`/`, Welcome)
	r.Post(`/status`, Status)
	r.Get(`/status`, Status)
	r.Get(`/last_status`, LastStatus)

	// Add the router action
	m.Action(r.Handle)
	m.Run()
}

func Welcome() string {
	return `Welcome to TeamDock, your TeamCity Flowdock connector`
}

func LastStatus() string {
	return fmt.Sprintf("Build: %+v, %+v", lastBuild, lastForm)
}

func Status(r *http.Request, w http.ResponseWriter) (result string) {
	currentBuild := Build{
		BuildStatusPrevious: r.FormValue("buildStatusPrevious"),
		BuildId:             r.FormValue("buildId"),
		TriggeredBy:         r.FormValue("triggeredBy"),
		BuildStatus:         r.FormValue("buildStatus"),
		AgentHostname:       r.FormValue("agentHostname"),
		BuildFullName:       r.FormValue("buildFullName"),
		BuildTypeId:         r.FormValue("buildTypeId"),
		Message:             r.FormValue("message"),
		Text:                r.FormValue("text"),
		NotifyType:          r.FormValue("notifyType"),
		AgentName:           r.FormValue("agentName"),
		BuildResult:         r.FormValue("buildResult"),
		BuildRunner:         r.FormValue("buildRunner"),
		ProjectId:           r.FormValue("projectId"),
		ProjectName:         r.FormValue("projectName"),
		AgentOs:             r.FormValue("agentOs"),
		BuildNumber:         r.FormValue("buildNumber"),
		BuildName:           r.FormValue("buildName"),
	}

	lastBuild = currentBuild
	lastForm = r.Form
	return "ok"
}
