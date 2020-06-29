package gp

import (
	"bytes"
	"github.com/cacarpenter/gopostal/cui"
	"github.com/cacarpenter/gopostal/postman"
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
	collections  []*postman.Collection
}

func (app *GoPostal) initLogging() {
	logfilename := os.TempDir() + "/gopostal.log"
	f, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// this means nothing else logs
	// defer f.Close()
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
	app.collections = make([]*postman.Collection, 1)

	return &app
}

func (app *GoPostal) ExecCurrentSelection() {
	pmColl := app.ui.SelectedCollection()
	if pmColl != nil && pmColl.Request != nil {
		response, err := app.CallRequest(pmColl.Request, app.logger.Writer())

		if err != nil {
			app.logger.Panicln(err)
		}

		for _, ev := range pmColl.Events {
			var buf bytes.Buffer
			for _, l := range ev.Script.Lines {
				buf.WriteString(l)
				buf.WriteString("\n")
			}
			app.RunJavaScript(buf.String(), *response)
		}
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

func (app *GoPostal) SetPostmanCollections(collections []*postman.Collection) {
	app.collections = collections
	app.ui.SetPostmanCollections(collections)
}

func (app *GoPostal) Run() {
	app.ui.Run()
}

func (app *GoPostal) Stop() {
	app.logger.Println("Bye")
	app.logFile.Close()
}
