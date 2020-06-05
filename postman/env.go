package postman

type Environment struct {
	Name      string
	variables map[string]string
}

func (e *Environment) Put(key, val string) {
	e.variables[key] = val
}

func (e *Environment) Get(key string) string {
	return e.variables[key]
}
