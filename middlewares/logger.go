package middlewares

import (
    "github.com/lovego/goa"
    "github.com/lovego/tracer"
    "time"
)

func Record(ctx *goa.Context){
    debug := ctx.URL.Query()["_debug"] != nil
    span := &tracer.Span{At: time.Now()}
    if debug {
        span.SetDebug(true)
    }
    var err error
    defer func() {
        panicErr := recover()
        if panicErr != nil && recoverFunc != nil {
            recoverFunc()
        }

        f := l.WithSpan(span)
        if fieldsFunc != nil {
            fieldsFunc(f)
        }

        if panicErr != nil {
            f.output(Recover, panicErr, f.data)
        } else if err != nil {
            f.output(Error, err, f.data)
        } else {
            f.output(Info, nil, f.data)
        }
    }()
}
