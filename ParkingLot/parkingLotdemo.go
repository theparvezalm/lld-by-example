package ParkingLot

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func ParkingLotDemo() {
	fmt.Println("=== Parking Lot System Demo ===\n")

	// Test 1: Build a lot and make sure setup works
	fmt.Println("Test 1: Creating parking lot...")
	var spots []*ParkingSpot
	spotType := VehicleTypesForService("BIKE")
	spot := NewParkingSpot("S1", spotType)
	spots = append(spots, spot)
	spotType = VehicleTypesForService("CAR")
	spot = NewParkingSpot("S2", spotType)
	spots = append(spots, spot)
	spotType = VehicleTypesForService("ALL")
	spot = NewParkingSpot("S3", spotType)
	spots = append(spots, spot)

	newParkingLot, err := NewParkingLot(spots)
	if err != nil {
		log.Fatalf("Failed to create parking lot: %v", err)
	}
	fmt.Printf("Parking lot created with capacity: %d\n", newParkingLot.GetTotalCapacity())
	fmt.Printf("Occupied spots: %d\n", newParkingLot.GetOccupiedSpotsCount())

	// Test 2: Try (and fail) to create an empty lot
	fmt.Println("\nTest 2: Testing empty spots validation...")
	emptySpots := []*ParkingSpot{}
	_, err = NewParkingLot(emptySpots)
	if err != nil {
		fmt.Printf("Correctly rejected empty parking lot: %v\n", err)
	}

	// Test 3: Let different vehicle types enter
	fmt.Println("\nTest 3: Multiple vehicle types entering...")
	ticket1, err := newParkingLot.Enter(MotorVehicleType)
	if err != nil {
		log.Fatalf("Failed to enter motor vehicle: %v", err)
	}
	fmt.Printf("Motor vehicle entered. Ticket: %s, Spot: %s\n", ticket1.Id, ticket1.SpotId)

	ticket2, err := newParkingLot.Enter(CAR)
	if err != nil {
		log.Fatalf("Failed to enter car: %v", err)
	}
	fmt.Printf("Car entered. Ticket: %s, Spot: %s\n", ticket2.Id, ticket2.SpotId)

	ticket3, err := newParkingLot.Enter(HeavyVehicle)
	if err != nil {
		log.Fatalf("Failed to enter heavy vehicle: %v", err)
	}
	fmt.Printf("Heavy vehicle entered. Ticket: %s, Spot: %s\n", ticket3.Id, ticket3.SpotId)

	// Test 4: See how many spots are left
	fmt.Println("\nTest 4: Checking available spots...")
	motorAvailable, _ := newParkingLot.GetAvailableSpotsCount(MotorVehicleType)
	carAvailable, _ := newParkingLot.GetAvailableSpotsCount(CAR)
	heavyAvailable, _ := newParkingLot.GetAvailableSpotsCount(HeavyVehicle)
	fmt.Printf("Available spots - Motor: %d, Car: %d, Heavy: %d\n", motorAvailable, carAvailable, heavyAvailable)
	fmt.Printf("Occupied spots: %d\n", newParkingLot.GetOccupiedSpotsCount())

	// Test 5: Expect a \"no spots left\" error
	fmt.Println("\nTest 5: Testing no available spots scenario...")
	_, err = newParkingLot.Enter(HeavyVehicle)
	if err != nil {
		fmt.Printf("Correctly handled no available spots: %v\n", err)
	}

	// Test 6: Exit vehicles and charge fees
	fmt.Println("\nTest 6: Exiting vehicles and calculating fees...")
	time.Sleep(100 * time.Millisecond) // Simulate some parking time
	fees1, err := newParkingLot.Exit(ticket1.Id)
	if err != nil {
		log.Fatalf("Failed to exit motor vehicle: %v", err)
	}
	fmt.Printf("Motor vehicle exited. Fee: %d\n", fees1)

	time.Sleep(100 * time.Millisecond)
	fees2, err := newParkingLot.Exit(ticket2.Id)
	if err != nil {
		log.Fatalf("Failed to exit car: %v", err)
	}
	fmt.Printf("Car exited. Fee: %d\n", fees2)

	time.Sleep(100 * time.Millisecond)
	fees3, err := newParkingLot.Exit(ticket3.Id)
	if err != nil {
		log.Fatalf("Failed to exit heavy vehicle: %v", err)
	}
	fmt.Printf("Heavy vehicle exited. Fee: %d\n", fees3)

	// Test 7: Break things with bad tickets
	fmt.Println("\nTest 7: Testing invalid ticket scenarios...")
	_, err = newParkingLot.Exit("invalid-ticket-id")
	if err != nil {
		fmt.Printf("Correctly rejected invalid ticket: %v\n", err)
	}
	_, err = newParkingLot.Exit("")
	if err != nil {
		fmt.Printf("Correctly rejected empty ticket ID: %v\n", err)
	}
	_, err = newParkingLot.Exit(ticket1.Id) // Try exiting with a ticket already used
	if err != nil {
		fmt.Printf("Correctly rejected already used ticket: %v\n", err)
	}

	// Test 8: Hammer the lot concurrently to see if locks hold up
	fmt.Println("\nTest 8: Testing concurrent access (thread safety)...")
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex
	concurrentTickets := make([]string, 0)

	// Spin up goroutines that all try to enter at once
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			vehicleType := MotorVehicleType
			if id%3 == 0 {
				vehicleType = CAR
			} else if id%3 == 1 {
				vehicleType = HeavyVehicle
			}
			ticket, err := newParkingLot.Enter(vehicleType)
			if err == nil {
				mu.Lock()
				successCount++
				concurrentTickets = append(concurrentTickets, ticket.Id)
				mu.Unlock()
				fmt.Printf("  Goroutine %d: Vehicle entered, Ticket: %s\n", id, ticket.Id)
			} else {
				fmt.Printf("  Goroutine %d: Failed to enter: %v\n", id, err)
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("Concurrent entries completed. Successful: %d\n", successCount)

	// Test 9: Let the concurrent cars leave
	fmt.Println("\nTest 9: Exiting concurrent vehicles...")
	for _, ticketID := range concurrentTickets {
		fee, err := newParkingLot.Exit(ticketID)
		if err == nil {
			fmt.Printf("  Exited ticket %s, Fee: %d\n", ticketID[:8], fee)
		}
	}

	// Test 10: Final sanity check
	fmt.Println("\nTest 10: Final state check...")
	finalAvailable, _ := newParkingLot.GetAvailableSpotsCount(MotorVehicleType)
	finalOccupied := newParkingLot.GetOccupiedSpotsCount()
	fmt.Printf("Final state - Available spots: %d, Occupied spots: %d, Total capacity: %d\n",
		finalAvailable, finalOccupied, newParkingLot.GetTotalCapacity())

	fmt.Println("\n=== All Tests Completed Successfully ===")
}
