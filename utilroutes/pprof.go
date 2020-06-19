package utilroutes

import (
	"bytes"
	"fmt"
	"html/template"
	hprof "net/http/pprof"
	"runtime"
	"runtime/pprof"
	"strconv"

	"github.com/lovego/goa"
)

// Pprof setup "runtime/pprof" profiles routes
func Pprof(router *goa.Router) {
	router.Group(`/_pprof`).Get(`/`, pprofIndex).Get(`/(.+)`, pprofGet)
}

var pprofIndexHtml []byte

func pprofIndex(c *goa.Context) {
	if pprofIndexHtml == nil {
		var tmpl = template.Must(template.New(``).Parse(`<html>
<head>
<title>pprof/</title>
</head>
<body>
pprof<br>
<br>
profiles:<br>
<table>
<tr> <td></td> <td><a href="/_pprof/profile">cpu profile</a></td> </tr>
<tr> <td></td> <td><a href="/_pprof/trace">trace</a></td> </tr>
{{range .}}
<tr>
<td align=right>{{.Count}}</td>
<td><a href="/_pprof/{{.Name}}?debug=1">{{.Name}}</a></td>
</tr>
{{end}}
</table>
<br>
<a href="/_pprof/goroutine?debug=2">full goroutine stack dump</a><br>
</body>
</html>
`))
		profiles := pprof.Profiles()
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, profiles); err != nil {
			c.Write([]byte(err.Error()))
		}
		pprofIndexHtml = buf.Bytes()
	}
	c.Write(pprofIndexHtml)
}

func pprofGet(c *goa.Context) {
	name := c.Param(0)
	switch name {
	case "profile":
		hprof.Profile(c.ResponseWriter, c.Request)
		return
	case "trace":
		hprof.Trace(c.ResponseWriter, c.Request)
		return
	}

	c.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	p := pprof.Lookup(name)
	if p == nil {
		c.WriteHeader(404)
		fmt.Fprintf(c, "Unknown profile: %s\n", name)
		return
	}
	if name == "heap" && c.FormValue("gc") != `` {
		runtime.GC()
	}
	debug, _ := strconv.Atoi(c.FormValue("debug"))
	p.WriteTo(c, debug)
}
