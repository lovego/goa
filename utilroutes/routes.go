package utilroutes

import (
	"github.com/lovego/goa"
)

func Setup(router *goa.Router) {
	router.Get(`/_alive`, func(c *goa.Context) {
		c.Write([]byte(`ok`))
	})

	Ps(router)
	Pprof(router)
}
