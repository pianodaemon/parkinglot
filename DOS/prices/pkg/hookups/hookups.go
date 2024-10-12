package hookups

import (
	"blaucorp.com/prices/internal/rest"
)

func CalisServer() {

	r := rest.GetEngine()
	rest.SetHandlers(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}
