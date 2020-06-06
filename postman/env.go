package postman

import "time"

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

func Read(filename string) (*Environment, error) {
	return nil, nil
}
