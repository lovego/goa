package middlewares

import (
    "net/http"
    "io"
    "os"
    "path/filepath"
    "encoding/json"
    "log"

    "github.com/lovego/goa"
    "github.com/lovego/tracer"
    "github.com/lovego/config"
    "github.com/lovego/fs"
    "github.com/lovego/slice"
    loggerPkg "github.com/lovego/logger"
    "io/ioutil"
)

var logBody = true
var logger = getLogger()
var logBodyMethods = []string{http.MethodPost, http.MethodDelete, http.MethodPut}

func getLogger() *loggerPkg.Logger {
    logger := loggerPkg.New(getLogWriter())
    logger.SetAlarm(config.Alarm())
    logger.SetMachineName()
    return logger
}

func getLogWriter() io.Writer {
    if config.DevMode() {
        return os.Stdout
    }
    file, err := fs.NewLogFile(filepath.Join(config.Root(), `log`, `app.log`))
    if err != nil {
        log.Fatal(err)
    }
    return file
}

func logFields(f *loggerPkg.Fields, ctx *goa.Context, debug bool) {
    req := ctx.Request
    res := ctx.Response
    f.With("host", req.Host)
    f.With("method", req.Method)
    f.With("path", req.URL.Path)
    f.With("rawQuery", req.URL.RawQuery)
    f.With("query", req.URL.Query())
    f.With("status", res.Status)
    f.With("reqBodySize", req.ContentLength)
    f.With("resBodySize", res.ContentLength)
    // 	f.With("proto", req.Proto)
    f.With("ip", ctx.ClientAddr())
    f.With("agent", req.UserAgent())
    f.With("refer", req.Referer())
    sess := GetSession(ctx)
    if sess != nil {
        f.With("session", sess)
    }
    reqBody, err1 := ioutil.ReadAll(req.Body)
    rspBody, err2 := ioutil.ReadAll(res.Body)
    if logBody && slice.ContainsString(logBodyMethods, req.Method) || debug {
        if err1 != nil{
            f.With("reqBody", tryUnmarshal(reqBody))
        }
        if err2 != nil{
            f.With("resBody", tryUnmarshal(rspBody))
        }
    }
}

func tryUnmarshal(b []byte) interface{} {
    var v map[string]interface{}
    err := json.Unmarshal(b, &v)
    if err == nil {
        return v
    } else {
        return string(b)
    }
}


func handleServerError(ctx *goa.Context) {
    ctx.WriteHeader(500)
    if ctx.Response.ContentLength <= 0 {
        ctx.Json(map[string]string{"code": "server-err", "message": "Fatal Server Error."})
    }
}

func Record(ctx *goa.Context){
    debug := ctx.URL.Query()["_debug"] != nil
    span := tracer.GetSpan(ctx.Context())
    if debug {
        span.SetDebug(true)
    }
    defer func() {
        panicErr := recover()
        if panicErr != nil {
            handleServerError(ctx)
        }
        f := logger.WithSpan(span)
        logFields(f, ctx, debug)

        if panicErr != nil {
            f.Recovery(panicErr)
        } else if err := ctx.GetError();err != nil {
            f.Error(err)
        } else {
            f.Info(nil)
        }
    }()
    ctx.Next()
    logger.Record(true, ctx.Next, nil, nil)

}
