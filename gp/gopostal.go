package gp

import (
	"fmt"
	"github.com/cacarpenter/gopostal/cui"
	"github.com/cacarpenter/gopostal/gpmodel"
	"io"
	"log"
	"os"
)

const GoPostalVersion = "0.1.0"

type GoPostal struct {
	ui        *cui.ConsoleUI
	session   *Session
	logger    *log.Logger
	logFile   *os.File
	varGroups []*gpmodel.VarGroup
	reqGroups []*gpmodel.RequestGroup
}

func (app *GoPostal) initLogging() {
	logfilename := os.TempDir() + "/gopostal.log"
	f, err := os.OpenFile(logfilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	app.logFile = f

	logger := log.New(f, "", log.Ldate|log.Ltime)
	logger.Println("GoPostal " + GoPostalVersion)
	app.logger = logger
}

func New() *GoPostal {
	app := GoPostal{}
	app.initLogging()
	app.ui = cui.NewConsoleUI(app.logger)
	app.ui.SetOnExec(app.ExecCurrentSelection)
	app.session = NewSession()
	app.varGroups = make([]*gpmodel.VarGroup, 1)
	app.reqGroups = make([]*gpmodel.RequestGroup, 1)

	return &app
}

func (app *GoPostal) ExecCurrentSelection(w io.Writer) {
	if app.ui.IsRequestSelected() {
		req := app.ui.SelectedRequest()
		w.Write([]byte(fmt.Sprintf("Calling %s\n", req.UrlPattern)))
		response, err := app.CallRequest(req, app.logger.Writer())

		if err != nil {
			app.logger.Println(err)
			w.Write([]byte(fmt.Sprintln(err)))
		} else {
			w.Write([]byte(fmt.Sprintf("Response:\n%s\n", *response)))
			app.RunJavaScript(req.PostScript.Text, *response)
		}
	}
}

func (app *GoPostal) SetVarGroups(varGroups []*gpmodel.VarGroup) {
	app.varGroups = varGroups
	for _, vg := range varGroups {
		app.logger.Println("Loading VarGroup Environment", vg.Name)
		app.session.Update(vg, true)
	}
	app.ui.UpdateVariables(app.session.variables)
}

func (app *GoPostal) SetRequestGroups(grps []*gpmodel.RequestGroup) {
	app.logger.Printf("Using %d request reqGroups\n", len(grps))
	for _, g := range grps {
		app.logger.Printf("Request Group %s\n", g.Name)
	}
	app.reqGroups = grps
	app.ui.SetRequestGroups(grps)
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
