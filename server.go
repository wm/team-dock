package main

import (
	"github.com/codegangsta/martini"
	"net/http"
	"fmt"
)

type Build struct {
	Branch        string   `json:"branch,omitempty"`
	Project       string   `json:"project,omitempty"`
	Result        string   `json:"result,omitempty"`
	Status        string   `json:"status,omitempty"`
	BuileNumber   string   `json:"buileNumber,omitempty"`
	TriggeredBy   string   `json:"triggeredBy,omitempty"`
	BuildTypeId   string   `json:"buildTypeId,omitempty"`
	BuildId       string   `json:"buildId,omitempty"`
}
	
var lastBuild = Build{}

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
	return fmt.Sprintf("Build: %+v", lastBuild)
}

func Status(r *http.Request, w http.ResponseWriter) (result string) {
	currentBuild := Build{
		Branch:      r.FormValue("branch"),
		Project:     r.FormValue("project"),
		Result:      r.FormValue("result"),
		Status:      r.FormValue("status"),
		BuileNumber: r.FormValue("number"),
		TriggeredBy: r.FormValue("by"),
		BuildTypeId: r.FormValue("typeId"),
		BuildId:     r.FormValue("buildId"),
	}

	result = fmt.Sprintf("Build: %+v", currentBuild)
	lastBuild = currentBuild
	return result
}
