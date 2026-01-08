package ParkingLot

import "time"

type Status int
type VehicleType int

const (
	MotorVehicleType VehicleType = iota
	CAR
	HeavyVehicle
)

const (
	Active Status = iota
	Closed
)

type Ticket struct {
	Id          string
	EntryTime   time.Time
	ExitTime    time.Time
	SpotId      string
	VehicleType VehicleType
	Status      Status
}

func NewTicket(id string, entryTime time.Time, spotId string, vehicleType VehicleType) Ticket {
	return Ticket{
		Id:          id,
		EntryTime:   entryTime,
		SpotId:      spotId,
		VehicleType: vehicleType,
		Status:      Active,
	}
}
