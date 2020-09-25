package gpmodel

// VarGroup represents a collection of key value pairs under a name, generally from a single source
type VarGroup struct {
	Name           string
	Variables      map[string]string
	SourceFilename string
}
