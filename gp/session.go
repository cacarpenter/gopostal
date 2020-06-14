package gp

import (
	"github.com/cacarpenter/gopostal/postman"
	"sync"
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

var singleton *Session
var once sync.Once

func CurrentSession() *Session {
	once.Do(func() {
		singleton = &Session{make(map[string]string)}
	})
	return singleton
}

func (s *Session) Update(env *postman.Environment, overwrite bool) {
	for _, v := range env.Values {
		if _, exists := s.variables[v.Key]; !exists || overwrite {
			s.variables[v.Key] = v.Value
		}
	}
}

func (s *Session) Vars() map[string]string {
	return s.variables
}
