package utilroutes

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/lovego/goa"
	"github.com/lovego/tracer"
)

var requests = requestsT{
	Instance: instanceName,
	Ps:       make(map[string]map[string]int),
}

type requestsT struct {
	Instance     string
	Ps           map[string]map[string]int
	sync.RWMutex `json:"-"`
}

func recordRequests(ctx *goa.Context) {
	request := ctx.Request
	var startTime time.Time
	if t := tracer.Get(ctx.Context()); t != nil {
		startTime = t.At
	} else {
		startTime = time.Now()
	}
	requests.Add(request.Method, ctx.URL.Path, startTime)
	defer requests.Remove(request.Method, request.URL.Path, startTime)

	ctx.Next()
}

func (ps *requestsT) Count() int {
	ps.RLock()
	defer ps.RUnlock()
	return len(ps.Ps)

}
func (ps *requestsT) ToJson() []byte {
	ps.RLock()
	defer ps.RUnlock()
	bytes, err := json.Marshal(ps)
	if err != nil {
		return []byte(fmt.Sprint(err))
	}
	return bytes
}

func (ps *requestsT) Add(method, path string, startTime time.Time) {
	ps.Lock()
	defer ps.Unlock()
	key := method + ` ` + path
	ts := startTime.Format(`2006-01-02T15:04:05Z0700`)
	if value, ok := ps.Ps[key]; ok {
		value[ts]++
	} else {
		ps.Ps[key] = map[string]int{ts: 1}
	}
}

func (ps *requestsT) Remove(method, path string, startTime time.Time) {
	ps.Lock()
	defer ps.Unlock()
	key := method + ` ` + path
	if value, ok := ps.Ps[key]; ok {
		if ts := startTime.Format(`2006-01-02T15:04:05Z0700`); value[ts] > 1 {
			value[ts]--
		} else if len(value) > 1 {
			delete(value, ts)
		} else {
			delete(ps.Ps, key)
		}
	}
}
