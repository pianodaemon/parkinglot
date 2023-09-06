package business

import (
	"errors"
	"fmt"
	"sync"

	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type (
	Car struct {
		uuid uuid.UUID
		year int
	}

	// The shelves to place car references
	CarPool map[string]*Car

	// Represents a parking lot type
	ParkingLot struct {
		logger    *logrus.Logger
		slots     CarPool
		ctrlMutex *sync.Mutex
	}

	// Value object of a Car
	// https://en.wikipedia.org/wiki/Value_object
	CarDTO struct {
		Identifier string `json:"id"`
		Year       int    `json:"year"`
	}
)

// Turns a uuid type into a trivial string
func strUUID(i uuid.UUID) string { return fmt.Sprintf("%s", i) }

// Produces a dto car instance from a real one
func renderDummyFromReal(car *Car) CarDTO {

	dto := CarDTO{}

	dto.Identifier = strUUID(car.uuid)
	dto.Year = car.year

	return dto
}

// Gets with reentrancy a dto car from its real one as per the car identifier
func (cc *ParkingLot) Display(carID string) (*CarDTO, error) {

	cc.ctrlMutex.Lock()
	j, found := cc.slots[carID]
	cc.ctrlMutex.Unlock()

	if found {
		jd := renderDummyFromReal(j)
		return &jd, nil
	}

	return nil, errors.New("Car not found")
}

// Gets with reentrancy all the dto cars from their real ones
func (cc *ParkingLot) DisplayAll() ([]CarDTO, error) {
	var dummies []CarDTO

	cc.ctrlMutex.Lock()
	for _, car := range cc.slots {
		dummies = append(dummies, renderDummyFromReal(car))
	}
	cc.ctrlMutex.Unlock()

	if dummies != nil {

		return dummies, nil
	}

	return nil, errors.New("There are not any car at the parking lot now")
}

// Attempts with reentrancy to destroy a car
func (cc *ParkingLot) Destroy(carID string) error {

	cc.logger.Printf("Attempting car destruction of %s", carID)

	cc.ctrlMutex.Lock()
	_, found := cc.slots[carID]
	defer cc.ctrlMutex.Unlock()

	if !found {
		return errors.New("Car not found")
	}

	delete(cc.slots, carID)
	cc.logger.Printf("Destroyed car %s", carID)
	return nil
}

// Spawns an newer instance of the parking lot
func NewParkingLot(logger *logrus.Logger) *ParkingLot {

	cc := &ParkingLot{
		logger:    logger,
		slots:     make(CarPool),
		ctrlMutex: &sync.Mutex{},
	}

	return cc
}
