package service

import (
	"blaucorp.com/prices/internal/rest"
)

// Engages the RESTful API
func Engage() {

	r := rest.GetEngine()
	rest.SetHandlers(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}
