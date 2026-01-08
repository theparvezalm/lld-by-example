package ParkingLot

type SpotType int

const (
	SPOT_MotorVehicleType SpotType = iota
	SPOT_CAR
	SPOT_HeavyVehicle
)

var MapOfVehicleToSpot = map[VehicleType]SpotType{
	MotorVehicleType: SPOT_MotorVehicleType,
	CAR:              SPOT_CAR,
	HeavyVehicle:     SPOT_HeavyVehicle,
}

func VehicleTypesForService(service string) []SpotType {
	switch service {
	case "BIKE":
		return []SpotType{SPOT_MotorVehicleType}
	case "CAR":
		return []SpotType{SPOT_MotorVehicleType, SPOT_CAR}
	case "ALL":
		return []SpotType{
			SPOT_MotorVehicleType,
			SPOT_CAR,
			SPOT_HeavyVehicle,
		}
	default:
		return nil
	}
}

type ParkingSpot struct {
	Id   string
	Type []SpotType
}

func NewParkingSpot(id string, types []SpotType) *ParkingSpot {
	return &ParkingSpot{
		Id:   id,
		Type: types,
	}
}
