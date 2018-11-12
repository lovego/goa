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
	"strconv"
	"strings"
	"time"

	"github.com/lovego/goa"
	"github.com/lovego/logger"
)

func ExampleLogger() {
	ts, buf := startTestServer()
	defer ts.Close()
	if _, err := http.Post(
		ts.URL+"?a=b&c=d", "application/json", strings.NewReader("[0,1,2]"),
	); err != nil {
		log.Panic(err)
	}
	printLog(buf)

	// Output:
	// {"code":"ok","message":"good"}
	//  31
	// agent: "Go-http-client/1.1"
	// host: "127.0.0.1:3000"
	// ip: "127.0.0.1"
	// level: "info"
	// method: "POST"
	// path: "/"
	// query: {"a":["b"],"c":["d"]}
	// rawQuery: "a=b&c=d"
	// refer: ""
	// reqBody: "[0,1,2]"
	// reqBodySize: 7
	// resBody: {"code":"ok","message":"good"}
	// resBodySize: 31
	// session: "session-data"
	// status: 200
}

func ExampleLogger_debug() {
	ts, buf := startTestServer()
	defer ts.Close()
	if _, err := http.Post(
		ts.URL+"?a=b&c=d&_debug", "application/json", strings.NewReader("[0,1,2]"),
	); err != nil {
		log.Panic(err)
	}
	printLog(buf)

	// Output:
	// {"code":"ok","message":"good"}
	//  31
	// agent: "Go-http-client/1.1"
	// host: "127.0.0.1:3000"
	// ip: "127.0.0.1"
	// level: "info"
	// method: "POST"
	// path: "/"
	// query: {"_debug":[""],"a":["b"],"c":["d"]}
	// rawQuery: "a=b&c=d&_debug"
	// refer: ""
	// reqBody: "[0,1,2]"
	// reqBodySize: 7
	// resBody: {"code":"ok","message":"good"}
	// resBodySize: 31
	// session: "session-data"
	// status: 200
}

func ExampleLogger_panic() {
	ts, buf := startTestServer()
	defer ts.Close()
	if _, err := http.Get(ts.URL + "/panic?_debug"); err != nil {
		log.Panic(err)
	}
	printLog(buf)

	// Output:
	// agent: "Go-http-client/1.1"
	// host: "127.0.0.1:3000"
	// ip: "127.0.0.1"
	// level: "recover"
	// method: "GET"
	// msg: "crash"
	// path: "/panic"
	// query: {"_debug":[""]}
	// rawQuery: "_debug"
	// refer: ""
	// reqBody: ""
	// reqBodySize: 0
	// resBody: {"code":"server-err","message":"Fatal Server Error."}
	// resBodySize: 54
	// status: 500
	// stack:  github.com/lovego/goa/middlewares.startTestServer.func2
}

func startTestServer() (*httptest.Server, *bytes.Buffer) {
	buf := bytes.NewBuffer(nil)
	router := goa.New()
	router.Use(NewLogger(logger.New(buf)).Record)
	router.Post("/", func(c *goa.Context) {
		c.Ok("good")
		c.Set("session", "session-data")
		resBody := c.ResponseBody()
		fmt.Println(string(resBody), len(resBody))
	})
	router.Get("/panic", func(c *goa.Context) {
		panic("crash")
	})

	ts := httptest.NewUnstartedServer(router)
	l, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Panic(err)
	}
	ts.Listener = l
	ts.Start()
	return ts, buf
}

func printLog(buf *bytes.Buffer) {
	var m = make(map[string]json.RawMessage)
	if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
		log.Panic(err)
	}

	at, duration, stack := m["at"], m["duration"], m["stack"]
	if _, err := time.Parse(time.RFC3339Nano, string(at[1:len(at)-1])); err != nil {
		log.Panic(err)
	}
	if d, err := strconv.ParseFloat(string(duration), 64); err != nil {
		log.Panic(err)
	} else if d <= 0 || d >= 1 {
		log.Panic(d)
	}
	delete(m, "at")
	delete(m, "duration")
	delete(m, "stack")

	for _, row := range sortMap(m) {
		fmt.Printf("%s: %s\n", row[0], row[1])
	}
	if len(stack) == 0 {
		return
	}
	var s string
	if err := json.Unmarshal(stack, &s); err != nil {
		log.Panic(err)
	}
	if i := strings.Index(s, "\n\t"); i > 0 {
		fmt.Println("stack: ", s[:i])
	} else {
		fmt.Println("stack: ", s)
	}
}

func sortMap(m map[string]json.RawMessage) (results [][2]interface{}) {
	for k, v := range m {
		results = append(results, [2]interface{}{k, v})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i][0].(string) < results[j][0].(string)
	})
	return
}
