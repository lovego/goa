package utilroutes

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/lovego/goa"
	"github.com/lovego/tracer"
)

var psData = psType{
	m: make(map[string]map[string]int),
}

type psType struct {
	sync.RWMutex
	m map[string]map[string]int
}

// Ps setup middleware and route to list all the requests in processing.
func Ps(router *goa.Router) {
	router.Use(processingList)
	router.Get(`/_ps`, func(c *goa.Context) {
		c.Write(psData.ToJson())
	})
}

// processingList list all the requests in processing.
func processingList(c *goa.Context) {
	request := c.Request
	var startTime time.Time
	if t := tracer.Get(c.Context()); t != nil {
		startTime = t.At
	} else {
		startTime = time.Now()
	}
	psData.Add(request.Method, c.URL.Path, startTime)
	defer psData.Remove(request.Method, request.URL.Path, startTime)

	c.Next()
}

func (ps *psType) ToJson() []byte {
	ps.RLock()
	defer ps.RUnlock()
	bytes, err := json.Marshal(ps.m)
	if err != nil {
		return []byte(fmt.Sprint(err))
	}
	return bytes
}

func (ps *psType) Add(method, path string, startTime time.Time) {
	ps.Lock()
	defer ps.Unlock()
	key := method + ` ` + path
	ts := startTime.Format(`2006-01-02T15:04:05Z0700`)
	if value, ok := ps.m[key]; ok {
		value[ts]++
	} else {
		ps.m[key] = map[string]int{ts: 1}
	}
}

func (ps *psType) Remove(method, path string, startTime time.Time) {
	ps.Lock()
	defer ps.Unlock()
	key := method + ` ` + path
	if value, ok := ps.m[key]; ok {
		if ts := startTime.Format(`2006-01-02T15:04:05Z0700`); value[ts] > 1 {
			value[ts]--
		} else if len(value) > 1 {
			delete(value, ts)
		} else {
			delete(ps.m, key)
		}
	}
}
