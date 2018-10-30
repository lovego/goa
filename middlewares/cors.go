package middlewares

import (
	"net/http"

	"github.com/lovego/goa"
)

type CORS struct {
	allow func(origin string) bool
}

func NewCORS(allow func(origin string) bool) CORS {
	return CORS{allow}
}

func (cors CORS) Check(c *goa.Context) {
	if origin := c.Request.Header.Get(`Origin`); origin != `` {
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
			header.Set(`Access-Control-Allow-Methods`, `GET, POST, PUT, DELETE`)
			header.Set(`Access-Control-Allow-Headers`, `X-Requested-With, Content-Type, withCredentials`)
			return
		}
	}

	c.Next()
}
