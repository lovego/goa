package utilroutes

import (
	"bytes"
	"runtime/pprof"
	"sort"
	"text/template"
)

func debugIndex() []byte {
	profiles := pprof.Profiles()
	sort.Slice(profiles, func(i, j int) bool {
		var orderI, orderJ = 99, 99
		for order, v := range []string{
			"heap", "allocs", "goroutine", "threadcreate", "block", "mutex",
		} {
			switch v {
			case profiles[i].Name():
				orderI = order
			case profiles[j].Name():
				orderJ = order
			}
		}
		return orderI < orderJ
	})

	var buf bytes.Buffer
	if err := debugIndexTmpl.Execute(&buf, struct {
		Instance string
		ReqCount int
		Profiles []*pprof.Profile
		Descs    map[string]string
	}{
		Instance: instanceName,
		ReqCount: requests.Count(),
		Profiles: profiles,
		Descs:    pprofDescs,
	}); err != nil {
		return []byte(err.Error())
	}
	return buf.Bytes()
}

var debugIndexTmpl = template.Must(template.New(``).Parse(`
<!doctype html>
<html>
	<head>
		<meta charset="UTF-8">
		<base target="_blank">
	<title>pprof</title>
		<style>
table { border-collapse: collapse; }
th,td { padding: 5px 10px; border: 1px dashed gray; }
		</style>
	</head>
	<body>
	<h3>{{ .Instance }} profiles</h3>
	<table>
		<tr>
			<th>Name</th>
			<th>Count</th>
			<th>Description</th>
		</tr>

		<tr>
			<td><a href=/_debug/reqs>requests</a></td>
			<td align=right>{{ .ReqCount }}</td>
			<td>Requests in processing.</td>
		</tr>
		<tr>
			<td><a href=/_debug/cpu>cpu</a></td>
			<td></td>
			<td> CPU profile. You can specify the duration in the "seconds" GET parameter. After you get the profile file, use the go tool pprof command to investigate the profile. </td>
		</tr>

		{{range .Profiles}}
		<tr>
			<td><a href="/_debug/{{.Name}}?debug=1">{{.Name}}</a></td>
			<td align=right>{{.Count}}</td>
			<td>{{ index $.Descs .Name }}</td>
		</tr>
		{{end}}

		<tr>
			<td><a href=/_debug/trace>trace</a></td>
			<td></td>
			<td> A trace of execution of the current program. You can specify the duration in the "seconds" GET parameter. After you get the trace file, use the go tool trace command to investigate the trace. </td>
		</tr>
	</table>
	</body>
</html>
`))

var pprofDescs = map[string]string{
	"heap":         "Heap memory profile. You can specify the gc GET parameter to run GC before taking the heap sample.",
	"allocs":       "Heap memory allocations, including all past.",
	"goroutine":    "Stack traces of all current goroutines.",
	"threadcreate": "Stack traces that led to the creation of new OS threads.",
	"block":        "Stack traces that led to blocking on synchronization primitives.",
	"mutex":        "Stack traces of holders of contended mutexes.",
}
