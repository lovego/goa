package middlewares

import (
    "github.com/lovego/goa"
    "github.com/lovego/session"
    "net/http"
)

const SessionKey = "goa_session"
const name = "sdf"
const secret = "dsf"

func CookieSession(ctx *goa.Context){
    cs := session.NewCookieSession(http.Cookie{
        Name: name,
        Path: "/",
    }, secret)
    ctx.Set(SessionKey, cs)
    ctx.Next()
}

func GetSession(ctx *goa.Context) *session.CookieSession{
    return ctx.Get(SessionKey).(*session.CookieSession)
}
