package gp

import (
	"github.com/cacarpenter/gopostal/postman"
)

type Session struct {
	variables map[string]string
}

func (s *Session) Put(key, val string) {
	s.variables[key] = val
}

func (s *Session) Get(key string) string {
	return s.variables[key]
}

func NewSession() *Session {
	return &Session{make(map[string]string)}
}

func (s *Session) Update(env *postman.Environment, overwrite bool) {
	for _, v := range env.Values {
		if _, exists := s.variables[v.Key]; !exists || overwrite {
			s.variables[v.Key] = v.Value
		}
	}
}

func (s *Session) Map() map[string]string {
	return s.variables
}

func (s *Session) Array() [][]string {
	// do some ordering?
	array := make([][]string, len(s.variables))
	i := 0
	for k, v := range s.variables {
		array[i] = make([]string, 2)
		array[i][0] = k
		array[i][1] = v
		i++
	}
	return array
}

func (s *Session) Names() []string {
	n := make([]string, len(s.variables))
	i := 0
	for k := range s.variables {
		n[i] = k
		i++
	}
	return n
}