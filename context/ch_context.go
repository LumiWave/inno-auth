package context

import (
	"sync"

	"github.com/gorilla/websocket"
)

const (
	WSChannel = "WSChannel"
)

var context *Context

var once sync.Once

func GetInstance() *Context {
	once.Do(func() {
		context = &Context{}
		context.data = make(map[string]interface{})
	})

	return context
}

type Context struct {
	data map[string]interface{}
}

func (o *Context) Put(key string, value interface{}) {
	o.data[key] = value
}

func (o *Context) Get(key string) (interface{}, bool) {
	val, exists := o.data[key]
	return val, exists
}

type ChData struct {
	Session     *websocket.Conn
	CommandType uint32

	Data     []byte
	Callback chan interface{}
}
