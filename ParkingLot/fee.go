package ParkingLot

import (
	"errors"
	"time"
)

const (
	MotorVehicleFeePerMinute = 10
	CarFeePerMinute          = 20
	HeavyVehicleFeePerMinute = 50
	MinimumChargeMinutes     = 60
)

type MotorVehicleFees struct{}

type CarFees struct{}

type HeavyVehicleFees struct{}

type FeeInterface interface {
	Calculate(entryTime, exitTime time.Time) (int64, error)
}

func (MotorVehicleFees) Calculate(entry, exit time.Time) (int64, error) {
	minutes := int64(exit.Sub(entry).Minutes())
	if minutes < MinimumChargeMinutes {
		minutes = MinimumChargeMinutes
	}
	return minutes * MotorVehicleFeePerMinute, nil
}

func (CarFees) Calculate(entry, exit time.Time) (int64, error) {
	minutes := int64(exit.Sub(entry).Minutes())
	if minutes < MinimumChargeMinutes {
		minutes = MinimumChargeMinutes
	}
	return minutes * CarFeePerMinute, nil
}

func (HeavyVehicleFees) Calculate(entry, exit time.Time) (int64, error) {
	minutes := int64(exit.Sub(entry).Minutes())
	if minutes < MinimumChargeMinutes {
		minutes = MinimumChargeMinutes
	}
	return minutes * HeavyVehicleFeePerMinute, nil
}

func GetFeeCalculator(vehicleType VehicleType) (FeeInterface, error) {
	switch vehicleType {
	case MotorVehicleType:
		return MotorVehicleFees{}, nil
	case CAR:
		return CarFees{}, nil
	case HeavyVehicle:
		return HeavyVehicleFees{}, nil
	default:
		return nil, errors.New("unknown vehicle type")
	}
}

func ComputeFees(entryTime, exitTime time.Time, vehicleType VehicleType) (int64, error) {
	calculator, err := GetFeeCalculator(vehicleType)
	if err != nil {
		return 0, err
	}
	if exitTime.Before(entryTime) {
		return 0, errors.New("exit time before entry time")
	}

	return calculator.Calculate(entryTime, exitTime)
}
