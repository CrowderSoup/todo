package store

import (
	"github.com/gin-gonic/gin"
	"github.com/go-session/gin-session"
	"github.com/go-session/session"
)

type Store interface {
	Get(string) (interface{}, bool)
	Set(string, interface{}) error
}

type store struct {
	session session.Store
}

func InitSession() gin.HandlerFunc {
	return ginsession.New()
}

func NewStore(c *gin.Context) Store {
	session := ginsession.FromContext(c)

	return &store{
		session,
	}
}

func (s store) Get(k string) (interface{}, bool) {
	return s.session.Get(k)
}

func (s store) Set(k string, v interface{}) error {
	s.session.Set(k, v)

	return s.session.Save()
}
