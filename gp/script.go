package gp

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"io"
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

func RunJavaScript(script, responseBody string, writer io.Writer) {
	vm := otto.New()

	if _, err := vm.Run(postmanBoostrapJs); err != nil {
		fmt.Fprintln(writer, "postman bootstrap returned error", err)
		return
	}
	vm.Set("responseBody", responseBody)

	// run the user script
	// TODO wrap in unsafe guards
	runVal, err := vm.Run(script)
	if err != nil {
		fmt.Fprintln(writer, "Script returned error", err)
		return
	}
	fmt.Fprintf(writer, "script runval %q\n", runVal)

	envVal, err := vm.Get("env")
	if err != nil {
		fmt.Fprintln(writer, "get env error", err)
		return
	}
	envObj := envVal.Object()
	for _, envKey := range envObj.Keys() {
		envVal, valErr := envObj.Get(envKey)
		if valErr != nil {
			continue
		}
		// TODO apply these to our session variables
		fmt.Fprintf(writer, "%s -> %s\n", envKey, envVal)
	}
}
