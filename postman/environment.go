package postman

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"
)

const POSTMAN_ENV_SUFFIX = ".postman_environment.json"

type Environment struct {
	Id                   string    `json:"id"`
	Name                 string    `json:"name"`
	Values               []EnvVal  `json:"values"`
	PostmanVariableScope string    `json:"_postman_variable_scope"`
	PostmanExportedAt    time.Time `json:"_postman_exported_at"`
	PostmanExportedUsing string    `json:"_postman_exported_using"`
}

type EnvVal struct {
	Key     string `json:"key"`
	Value   string `json:"value"`
	Enabled bool   `json:"enabled"`
}

func ParseEnv(filename string) (*Environment, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var e Environment
	err = json.Unmarshal(data, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func IsEnvironmentFile(filename string) bool {
	return strings.HasSuffix(filename, POSTMAN_ENV_SUFFIX)
}

