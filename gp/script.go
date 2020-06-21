package gp

import (
	"fmt"
	"github.com/cacarpenter/gopostal/postman"
	"github.com/robertkrimen/otto"
)

type pm struct {
	environment postman.Environment
}

func RunJavaScript(script, responseBody string) {
	fmt.Println(script)
	vm := otto.New()

	var pmjs pm
	pmjs.environment = postman.Environment{}

	vm.Set("pm", pmjs)
	vm.Set("responseBody", responseBody)

	runVal, err := vm.Run(script)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(runVal)
}
