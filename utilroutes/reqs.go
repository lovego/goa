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
	Reqs:     make(map[string]map[string]int),
}

type requestsT struct {
	Instance     string
	Reqs         map[string]map[string]int
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

func (reqs *requestsT) Count() int {
	reqs.RLock()
	defer reqs.RUnlock()
	return len(reqs.Reqs)

}
func (reqs *requestsT) ToJson() []byte {
	reqs.RLock()
	defer reqs.RUnlock()
	bytes, err := json.Marshal(reqs)
	if err != nil {
		return []byte(fmt.Sprint(err))
	}
	return bytes
}

func (reqs *requestsT) Add(method, path string, startTime time.Time) {
	reqs.Lock()
	defer reqs.Unlock()
	key := method + ` ` + path
	ts := startTime.Format(`2006-01-02T15:04:05Z0700`)
	if value, ok := reqs.Reqs[key]; ok {
		value[ts]++
	} else {
		reqs.Reqs[key] = map[string]int{ts: 1}
	}
}

func (reqs *requestsT) Remove(method, path string, startTime time.Time) {
	reqs.Lock()
	defer reqs.Unlock()
	key := method + ` ` + path
	if value, ok := reqs.Reqs[key]; ok {
		if ts := startTime.Format(`2006-01-02T15:04:05Z0700`); value[ts] > 1 {
			value[ts]--
		} else if len(value) > 1 {
			delete(value, ts)
		} else {
			delete(reqs.Reqs, key)
		}
	}
}
