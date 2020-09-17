package gp

import (
	"fmt"
	"github.com/cacarpenter/gopostal/cui"
	"github.com/cacarpenter/gopostal/gpmodel"
	"github.com/cacarpenter/gopostal/postman"
	"io"
	"log"
	"os"
)

const GO_POSTAL_VERSION = "0.1.0"

type GoPostal struct {
	ui           *cui.ConsoleUI
	session      *Session
	logger       *log.Logger
	logFile      *os.File
	environments []*postman.Environment
	groups       []*gpmodel.Group
}

func (app *GoPostal) initLogging() {
	logfilename := os.TempDir() + "/gopostal.log"
	f, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	app.logFile = f

	logger := log.New(f, "", log.Ldate|log.Ltime)
	logger.Println("GoPostal " + GO_POSTAL_VERSION)
	app.logger = logger
}

func New() *GoPostal {
	app := GoPostal{}
	app.initLogging()
	app.ui = cui.NewConsoleUI(app.logger)
	app.ui.SetOnExec(app.ExecCurrentSelection)
	app.session = NewSession()
	app.environments = make([]*postman.Environment, 1)
	app.groups = make([]*gpmodel.Group, 1)

	return &app
}

func (app *GoPostal) ExecCurrentSelection(w io.Writer) {
	if app.ui.IsRequestSelected() {
		req := app.ui.SelectedRequest()
		w.Write([]byte(fmt.Sprintf("Calling %s\n", req.UrlPattern)))
		response, err := app.CallRequest(req, app.logger.Writer())
		w.Write([]byte(fmt.Sprintf("Response is %s\n", response)))

		if err != nil {
			app.logger.Println(err)
			w.Write([]byte(fmt.Sprintln(err)))
		} else {
			w.Write([]byte(fmt.Sprintf("%q\n", *response)))
			app.RunJavaScript(req.PostScript.Text, *response)
		}
	} else {
		//		app.ui
	}
}

func (app *GoPostal) SetPostmanEnvironments(environments []*postman.Environment) {
	for _, pmEnv := range environments {
		app.logger.Println("Loading Postman Environment", pmEnv.Name)
		app.session.Update(pmEnv, true)
	}
	app.environments = environments
	app.ui.UpdateVariables(app.session.variables)
}

func (app *GoPostal) SetGroups(grps []*gpmodel.Group) {
	app.logger.Printf("Using %d groups\n", len(grps))
	for _, g := range grps {
		app.logger.Printf("Collection %s\n", g.Name)
	}
	app.groups = grps
	app.ui.SetGroups(grps)
}

func (app *GoPostal) Run() {
	app.ui.Run()
}

func (app *GoPostal) Stop() {
	app.logger.Println("Bye")
	app.logFile.Close()
}

func (app *GoPostal) UpdateSession(key, val string) {
	app.ui.UpdateVariable(key, val)
	app.session.Put(key, val)
}
