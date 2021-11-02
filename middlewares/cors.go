package middlewares

import (
	"net/http"

	"github.com/lovego/goa"
)

// crosss origin resource share
type CORS struct {
	allow     func(origin string) bool
	SetHeader func(http.Header)
}

func NewCORS(allow func(origin string) bool) CORS {
	return CORS{allow: allow}
}

func (cors CORS) Check(c *goa.Context) {
	if c.Request.Header.Get(`Sec-Fetch-Site`) != "same-origin" {
		if origin := c.Request.Header.Get(`Origin`); origin != `` && origin != c.Origin() {
			if !cors.allow(origin) {
				c.WriteHeader(http.StatusForbidden)
				c.Write([]byte(`origin not allowed.`))
				return
			}
			header := c.ResponseWriter.Header()
			header.Set(`Access-Control-Allow-Origin`, origin)
			header.Set(`Access-Control-Allow-Credentials`, `true`)
			header.Set(`Vary`, `Accept-Encoding, Origin`)

			if c.Request.Method == `OPTIONS` { // preflight request
				header.Set(`Access-Control-Max-Age`, `86400`)
				header.Set(`Access-Control-Allow-Methods`, `GET, POST, PUT, DELETE, PATCH`)
				header.Set(`Access-Control-Allow-Headers`, `X-Requested-With, Content-Type, withCredentials`)
				if cors.SetHeader != nil {
					cors.SetHeader(header)
				}
				return
			}
		}
	}

	c.Next()
}
