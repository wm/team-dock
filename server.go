package main

import (
	"github.com/codegangsta/martini"
	"github.com/wm/go-flowdock/flowdock"
	"net/http"
	"net/url"
	"fmt"
	"syscall"
	"text/template"
	"bytes"
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

//TODO: use const for uri's etc.
var flowClient *flowdock.Client //TODO pass this into the context or martini
var lastActionErr error
var lastBuild = Build{}
var lastForm = url.Values{} //map[string]string{}
var lastAction = "No action yet taken"

func main() {
	flowClient = flowdock.NewClient(nil)
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
	return fmt.Sprintf("Build: %+v, %+v, %+v, %+v", lastBuild, lastForm, lastAction, lastActionErr)
}

func Status(r *http.Request, w http.ResponseWriter) (result string) {
	build := encodeBuild(r)
	apiToken, ok := syscall.Getenv("FLOW_API_TOKEN")
	if ok {
		lastAction = "sending to Flow"
		lastActionErr = sendBuildToFlow(flowClient, build, apiToken)
	} else {
		lastActionErr = nil
		lastAction = "FLOW_API_TOKEN is not set"
	}

	lastForm = r.Form
	lastBuild = *build
	return "ok"
}

// TODO: do it similar to json Encode method in the future
func encodeBuild(r *http.Request) *Build {
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

	return &currentBuild
}

func sendBuildToFlow(client *flowdock.Client, build *Build, flowApiToken string) error {
	var fromAddress string

	switch build.BuildResult {
	case "success":  fromAddress = "build+ok@flowdock.com"
	case "running":  fromAddress = "build+started@flowdock.com"
	case "canceled": fromAddress = "build+canceled@flowdock.com"
	default:         fromAddress = "build+fail@flowdock.com"
	}

	body := statusBody(build)
	opt := &flowdock.InboxCreateOptions{
		Source:       "go-flowdock",
		FromAddress:  fromAddress,
		Subject:     fmt.Sprintf("%v build %v - %v", build.ProjectName, build.BuildNumber, build.BuildResult),
		Tags:        []string{build.BuildResult, "CI", build.BuildNumber, build.ProjectName},
		Project:     build.ProjectName,
		FromName:    "TeamCity CI",
		Content:      body,
	}

	_, _, err := client.Inbox.Create(flowApiToken, opt)

	return err
}

func statusBody(build *Build) string {
	var body bytes.Buffer

	bodyTmpl, err  := template.New("body").Parse(`
<ul>
	<li>
	<code><a href="https://github.com/IoraHealth/{{.ProjectName}}">IoraHealth/{{.ProjectName}}</a></code>
	build #{{.BuildNumber}}: {{.BuildResult}}!
	</li>
	<li>
		Branch: <code>{{.BuildName}}</code>
	</li>
	<li>
	<a href="http://nest.icisapp.com/viewLog.html?buildId={{.BuildId}}&tab=buildLog&buildTypeId={{.BuildTypeId}}">Build details</a>
	</li>
	<li>
    {{.Message}}
	</li>
</ul>
`)

	if err != nil { panic(err) }
	err = bodyTmpl.Execute(&body, build)
	if err != nil { panic(err) }

    return body.String()
}
