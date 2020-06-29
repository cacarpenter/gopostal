package gp

import (
	"github.com/robertkrimen/otto"
	"log"
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

func RunJavaScript(script, responseBody string) {
	vm := otto.New()

	if _, err := vm.Run(postmanBoostrapJs); err != nil {
		log.Println("postman bootstrap returned error", err)
		return
	}
	vm.Set("responseBody", responseBody)

	// run the user script
	// TODO wrap in unsafe guards
	runVal, err := vm.Run(script)
	if err != nil {
		log.Println("Script returned error", err)
		return
	}
	log.Printf("script runval %q\n", runVal)

	envVal, err := vm.Get("env")
	if err != nil {
		log.Println( "get env error", err)
		return
	}
	envObj := envVal.Object()
	log.Println(envObj)

	// TODO move this logic, dont access the session here
	/*
	session := CurrentSession()
	for _, envKey := range envObj.Keys() {
		envVal, valErr := envObj.Get(envKey)
		if valErr != nil {
			continue
		}
		log.Printf("%s -> %s\n", envKey, envVal)
		session.Put(envKey, envVal.String())
	}
	 */
	// TODO do something to trigger rerender of variables view
}
