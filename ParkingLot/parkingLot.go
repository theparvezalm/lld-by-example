package ParkingLot

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ParkingLot struct {
	mu              sync.RWMutex
	Spots           []*ParkingSpot
	ActiveTickets   map[string]Ticket
	occupiedSpotIds map[string]bool
}

func NewParkingLot(spots []*ParkingSpot) (*ParkingLot, error) {
	if len(spots) == 0 {
		return nil, errors.New("parking lot must have at least one spot")
	}
	return &ParkingLot{
		Spots:           spots,
		ActiveTickets:   make(map[string]Ticket),
		occupiedSpotIds: make(map[string]bool),
	}, nil
}

func (p *ParkingLot) Enter(vehicleType VehicleType) (Ticket, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	spot, err := p.findAvailableSpots(vehicleType)
	if err != nil {
		return Ticket{}, err
	}
	if spot == nil {
		return Ticket{}, errors.New("no available spots for vehicle type")
	}
	// Validate that spot exists in the parking lot
	spotExists := false
	for _, s := range p.Spots {
		if s.Id == spot.Id {
			spotExists = true
			break
		}
	}
	if !spotExists {
		return Ticket{}, errors.New("spot not found in parking lot")
	}
	p.occupiedSpotIds[spot.Id] = true
	ticketID := uuid.New().String()
	entryTime := time.Now()
	ticket := NewTicket(ticketID, entryTime, spot.Id, vehicleType)
	p.ActiveTickets[ticketID] = ticket
	return ticket, nil
}

func (p *ParkingLot) Exit(ticketId string) (int64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ticketId == "" {
		return 0, errors.New("invalid ticket ID")
	}
	ticket, exists := p.ActiveTickets[ticketId]
	if !exists {
		return 0, errors.New("ticket not found or already used")
	}
	if ticket.Status == Closed {
		return 0, errors.New("ticket is closed")
	}
	exitTime := time.Now()
	fee, err := ComputeFees(ticket.EntryTime, exitTime, ticket.VehicleType)
	if err != nil {
		return 0, err
	}
	ticket.Status = Closed
	ticket.ExitTime = exitTime
	p.ActiveTickets[ticketId] = ticket
	delete(p.occupiedSpotIds, ticket.SpotId)
	delete(p.ActiveTickets, ticketId)

	return fee, err
}

func (p *ParkingLot) findAvailableSpots(vehicleType VehicleType) (*ParkingSpot, error) {
	spotType, ok := MapOfVehicleToSpot[vehicleType]
	if !ok {
		return nil, errors.New("invalid vehicle type")
	}
	for _, spot := range p.Spots {
		if p.occupiedSpotIds[spot.Id] {
			continue
		}
		if supports(spot.Type, spotType) {
			return spot, nil
		}
	}
	return nil, nil
}

func (p *ParkingLot) GetAvailableSpotsCount(vehicleType VehicleType) (int, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	spotType, ok := MapOfVehicleToSpot[vehicleType]
	if !ok {
		return 0, errors.New("invalid vehicle type")
	}
	count := 0
	for _, spot := range p.Spots {
		if !p.occupiedSpotIds[spot.Id] && supports(spot.Type, spotType) {
			count++
		}
	}
	return count, nil
}

func (p *ParkingLot) GetTotalCapacity() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.Spots)
}

func (p *ParkingLot) GetOccupiedSpotsCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.occupiedSpotIds)
}

func supports(types []SpotType, required SpotType) bool {
	for _, t := range types {
		if t == required {
			return true
		}
	}
	return false
}
