package postman

import (
	"encoding/json"
	"github.com/cacarpenter/gopostal/gpmodel"
	"io/ioutil"
	"strings"
	"time"
)

const PostmanEnvSuffix = ".postman_environment.json"

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

func ParseEnvironment(filename string) (*gpmodel.VarGroup, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var pmEnv Environment
	err = json.Unmarshal(data, &pmEnv)
	if err != nil {
		return nil, err
	}

	vg := gpmodel.VarGroup{}
	vg.SourceFilename = filename
	vars := make(map[string]string, len(pmEnv.Values))
	for _, pmEnvVar := range pmEnv.Values {
		vars[pmEnvVar.Key] = pmEnvVar.Key
	}

	vg.Variables = vars
	return &vg, nil
}

func IsEnvironmentFile(filename string) bool {
	return strings.HasSuffix(filename, PostmanEnvSuffix)
}
