package service

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	co "immortalcrab.com/parkinglot/internal/controllers"
	"immortalcrab.com/parkinglot/internal/rsapi"
	ve "immortalcrab.com/parkinglot/pkg/business"
)

const carIDMask string = "[[:alnum:]\\-]+"

var apiSettings rsapi.RestAPISettings

func init() {

	envconfig.Process("rsapi", &apiSettings)
}

// Engages the RESTful API
func Engage(logger *logrus.Logger) (merr error) {

	defer func() {

		if r := recover(); r != nil {
			merr = r.(error)
		}
	}()

	cc := ve.NewParkingLot(logger)

	/* The connection of both components occurs through
	   the router glue and its adaptive functions */
	glue := func(api *rsapi.RestAPI) *mux.Router {

		router := mux.NewRouter()

		v1 := router.PathPrefix("/v1").Subrouter()
		mgmt := v1.PathPrefix("/crud").Subrouter()
		mgmt.HandleFunc("/cars", co.ListCars(cc.DisplayAll)).Methods("GET")
		mgmt.HandleFunc(fmt.Sprintf("/{car_id:%s}", carIDMask), co.ListCar(cc.Display)).Methods("GET")
		mgmt.HandleFunc(fmt.Sprintf("/{car_id:%s}/delete", carIDMask), co.DeleteCar(cc.Destroy)).Methods("PUT")
		return router
	}

	api := rsapi.NewRestAPI(logger, &apiSettings, glue)

	api.PowerOn()

	return nil
}
