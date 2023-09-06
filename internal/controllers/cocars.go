package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/jsonapi"
	ve "immortalcrab.com/parkinglot/pkg/business"
)

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
