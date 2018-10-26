package middlewares

import (
	"encoding/json"
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

// Ps records the processing requests.
func Ps(ctx *goa.Context) {
	request := ctx.Request
	var startTime time.Time
	if span := tracer.GetSpan(ctx.Context()); span != nil {
		startTime = span.At
	} else {
		startTime = time.Now()
	}
	psData.Add(request.Method, ctx.URL.Path, startTime)
	defer psData.Remove(request.Method, request.URL.Path, startTime)

	ctx.Next()
}

func (ps *psType) ToJson() []byte {
	ps.RLock()
	defer ps.RUnlock()
	bytes, err := json.Marshal(ps.m)
	if err != nil {
		panic(err)
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
