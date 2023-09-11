package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestEmptyPersitance(t *testing.T) {

	logger, _ := test.NewNullLogger()
	cc := ve.NewParkingLot(logger)

	req, _ := http.NewRequest("GET", "/v1/cruds/cars/list", nil)
	response := executeRequest(configureRouter(cc), req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}
