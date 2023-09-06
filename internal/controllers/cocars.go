package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	ve "immortalcrab.com/parkinglot/pkg/business"
)

// Deletes a Car
func DeleteCar(destructor func(string) error) func(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		Code  int    `json:"code"`
		CarID string `json:"car_id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		carID := vars["car_id"]

		w.Header().Set("Content-Type", "application/json")

		if err := destructor(carID); err != nil {

			w.WriteHeader(http.StatusNotFound)

			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Code:   strconv.Itoa(int(EndPointFailedDeletion)),
				Title:  "Failed deletion",
				Detail: err.Error(),
			}})

			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(Response{
			Code:  int(Success),
			CarID: carID,
		})

	}
}

// Displays a choosen car
func ListCar(displayer func(string) (*ve.CarDTO, error)) func(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		Info ve.CarDTO `json:"car"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		carID := vars["car_id"]

		dto, err := displayer(carID)

		w.Header().Set("Content-Type", "application/json")

		if err != nil {

			w.WriteHeader(http.StatusNotFound)

			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Code:   strconv.Itoa(int(EndPointCarNotFound)),
				Title:  "Car not found",
				Detail: err.Error(),
			}})

			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(Response{
			Info: *dto,
		})
	}
}

// Displays all the existing cars
func ListCars(pullInfo func() ([]ve.CarDTO, error)) func(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		Cars []ve.CarDTO `json:"cars"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		dtos, err := pullInfo()

		w.Header().Set("Content-Type", "application/json")

		if err != nil {

			w.WriteHeader(http.StatusBadRequest)

			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Code:   strconv.Itoa(int(EndPointNoCarsYet)),
				Title:  "No cars yet",
				Detail: err.Error(),
			}})

			return
		}

		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(Response{
			Cars: dtos,
		})
	}
}
