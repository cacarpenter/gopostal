package gp

import (
	"github.com/robertkrimen/otto"
)


// bootstrap postman variable assumed by scripts
const postmanBoostrapJs = `
var env = {}
var responseBody = ""

var pm = {};
pm.environment = {};
pm.environment.set = function(name, val) { env[name] = val };
pm.response = {};
pm.response.json = function() { return JSON.parse(responseBody); };
pm.test = function(name, testfunc) { testfunc(); };
`

func (app *GoPostal) RunJavaScript(script, responseBody string) {
	vm := otto.New()

	if _, err := vm.Run(postmanBoostrapJs); err != nil {
		app.logger.Println("postman bootstrap returned error", err)
		return
	}
	vm.Set("responseBody", responseBody)

	// run the user script
	// TODO wrap in unsafe guards
	runVal, err := vm.Run(script)
	if err != nil {
		app.logger.Println("Script returned error", err)
		return
	}
	app.logger.Printf("script runval %q\n", runVal)

	envVal, err := vm.Get("env")
	if err != nil {
		app.logger.Println( "get env error", err)
		return
	}
	envObj := envVal.Object()
	app.logger.Println(envObj)

	for _, envKey := range envObj.Keys() {
		envVal, valErr := envObj.Get(envKey)
		if valErr != nil {
			continue
		}
		app.logger.Printf("%s -> %s\n", envKey, envVal)
		app.ui.UpdateVariable(envKey, envVal.String())
	}
	// TODO do something to trigger rerender of variables view
}
