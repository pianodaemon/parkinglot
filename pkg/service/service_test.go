package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	co "immortalcrab.com/parkinglot/internal/controllers"
	ve "immortalcrab.com/parkinglot/pkg/business"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus/hooks/test"
)

func executeRequest(router *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func glue(cc *ve.ParkingLot) *mux.Router {

	router := mux.NewRouter()
	v1 := router.PathPrefix("/v1").Subrouter()
	cruds := v1.PathPrefix("/cruds").Subrouter()

	mgmt := cruds.PathPrefix("/cars").Subrouter()
	mgmt.HandleFunc("/list", co.ListCars(cc.DisplayAll)).Methods("GET")
	mgmt.HandleFunc(fmt.Sprintf("/{car_id:%s}", carIDMask), co.ListCar(cc.Display)).Methods("GET")
	mgmt.HandleFunc(fmt.Sprintf("/{car_id:%s}/delete", carIDMask), co.DeleteCar(cc.Destroy)).Methods("DELETE")
	mgmt.HandleFunc("/create", co.CreateCar(cc.Place)).Methods("POST")
	return router
}

func TestEmptyPersitance(t *testing.T) {

	logger, _ := test.NewNullLogger()
	cc := ve.NewParkingLot(logger)

	req, _ := http.NewRequest("GET", "/v1/cruds/cars/list", nil)
	response := executeRequest(glue(cc), req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}
