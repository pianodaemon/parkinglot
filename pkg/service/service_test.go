package service

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func newParkingLotWithNullLogging() *ve.ParkingLot {
	logger, _ := test.NewNullLogger()
	return ve.NewParkingLot(logger)
}

func TestEmptyPersitance(t *testing.T) {

	cc := newParkingLotWithNullLogging()

	req, _ := http.NewRequest("GET", "/v1/cruds/cars/list", nil)
	response := executeRequest(configureRouter(cc), req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestCarCreationAndDeletion(t *testing.T) {

	cc := newParkingLotWithNullLogging()
	router := configureRouter(cc)

	var jsonStr = []byte(`{ "year":2020 }`)
	req, _ := http.NewRequest("POST", "/v1/cruds/cars/create", bytes.NewBuffer(jsonStr))
	response := executeRequest(router, req)
	checkResponseCode(t, http.StatusAccepted, response.Code)
	m := struct {
		CarID string `json:"car_id"`
	}{}

	json.Unmarshal(response.Body.Bytes(), &m)
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/v1/cruds/cars/%s/delete", m.CarID), nil)
	response = executeRequest(router, req)
	checkResponseCode(t, http.StatusOK, response.Code)
}
