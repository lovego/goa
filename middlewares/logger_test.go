package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"sort"
	"time"

	"github.com/lovego/goa"
	"github.com/lovego/logger"
)

func ExampleRecord() {
	buf := bytes.NewBuffer(nil)
	router := goa.New()
	router.Use(NewLogger(logger.New(buf)).Record)
	router.Get("/", func(c *goa.Context) {
		c.Ok("ok")
	})
	ts := httptest.NewUnstartedServer(router)
	l, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Panic(err)
	}
	ts.Listener = l
	ts.Start()
	defer ts.Close()

	if _, err := http.Get(ts.URL + "?a=b&c=d"); err != nil {
		log.Panic(err)
	}
	var m = make(map[string]interface{})
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		log.Panic(err)
	}

	at, duration := m["at"], m["duration"]
	if _, err := time.Parse(time.RFC3339Nano, at.(string)); err != nil {
		log.Panic(err)
	}
	if d := duration.(float64); d <= 0 || d >= 1 {
		log.Panic(d)
	}
	delete(m, "at")
	delete(m, "duration")

	for _, row := range sortMap(m) {
		if v, ok := row[1].(string); !ok || v != "" {
			fmt.Printf("%s: %v\n", row[0], row[1])
		}
	}

	// Output:
	// agent: Go-http-client/1.1
	// host: 127.0.0.1:3000
	// ip: 127.0.0.1
	// level: info
	// method: GET
	// path: /
	// query: map[a:[b] c:[d]]
	// rawQuery: a=b&c=d
	// reqBodySize: 0
	// resBodySize: 28
	// status: 200
}

func sortMap(m map[string]interface{}) (results [][2]interface{}) {
	for k, v := range m {
		results = append(results, [2]interface{}{k, v})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i][0].(string) < results[j][0].(string)
	})
	return
}
