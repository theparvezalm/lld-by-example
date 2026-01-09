What is a Parking Lot System?
A parking lot system manages vehicle parking across multiple spots. When a vehicle enters, the system assigns an available spot matching the vehicle type and issues a ticket. When the vehicle exits, the system calculates the parking fee based on time spent and frees up the spot for the next customer.


Clarifying Questions
The system assigns an available spot matching the vehicle type and issues a ticket.
Vehical types. Motorcycles, regular cars, and large vehicles like SUVs or vans.
When a vehical enters they get a ticket with a unique ID. They’ll need that ticket to exit.
Pricing- Keep it simple. Hourly rate, same for all vehicles. Round up to the nearest hour and they pay when they leave.
Reject entry if there’s no compatible spot available. For exit, return an error if the ticket is invalid or already used.
Focus on the core logic. Spot assignment, ticket management, fee calculation. Skip the physical hardware, payment systems, and UI.
Requirements:
1. System supports three vehicle types: Motorcycle, Car, Large Vehicle
2. When a vehicle enters, system automatically assigns an available compatible spot
3. System issues a ticket at entry.
4. When a vehicle exits, user provides ticket ID
   — System validates the ticket
   — Calculates fee based on time spent (hourly, rounded up)
   — Frees the spot for next use
5. Pricing is hourly with same rate for all vehicles
6. System rejects entry if no compatible spot is available
7. System rejects exit if ticket is invalid or already used

Out of scope:
- Payment processing
- Physical gate hardware
- Security cameras or monitoring
- UI/display systems
- Reservations or pre-booking

Core Entities and Relationships
Vehical
At first glance, it feels natural to create a Vehicle class because the system deals with parking vehicles. However, if you think a bit deeper, the vehicle itself is not something our system owns or manages. We are not tracking its lifecycle, behavior, or internal state.For our use case, the only information we actually need is the vehicle’s category — such as motorcycle, car, or large — so that we can assign it to an appropriate parking spot. This is just a simple classification, not a full-fledged domain object.Because of that, modeling a Vehicle as a class would be unnecessary. An enum is sufficient and keeps the design cleaner and simpler.
ParkingSpot
- id
- type (This represents the fixed nature of the parking spot.What can this spot accept?motorcycle, car, large)
- isOccupied
  ParkingSpotHistory
- id
- parkingSpotId
- vehicalType
- occupiedAt
- vacatedAt
- status (e.g., OCCUPIED, VACATED, CANCELLED).
  Ticket
- id
- entryTime
- exitTime
- VehicalTypE
- status (ACTIVE, CLOSED)
- spotId
  ParkingLot
  The system needs a central component to manage all operations. When a vehicle comes in, there must be something responsible for locating a free spot, issuing a parking ticket, and updating the spot’s status to occupied. Similarly, when a vehicle leaves, that same component should verify the ticket, compute the parking charges, and release the spot.This role is handled by the ParkingLot. It acts as the main entry point of the system and coordinates interactions between different parts, ensuring the overall flow works correctly.
  The ParkingLot is responsible for managing all ParkingSpots within it. When a vehicle enters, the ParkingLot generates a Ticket. It also maintains a mapping of currently active tickets, allowing it to quickly retrieve a ticket using its ID when a vehicle exits.

Class Design
Important — Go does NOT have classes.Instead, Go uses:
struct → similar to a class (data)
methods → functions attached to a struct
Together, they behave like classes and objects.

For each class, we’ll ask two questions:
1. What does this class need to remember to enforce the requirements (its state)?
2. What operations does this class need to support (its methods)?

ParkingLot
ParkingLot is the controller. Everything flows through it. Let’s derive its state from requirements:
ParkingLot must track:

System automatically assigns an available compatible spot:
- All parking spots in the lot
- Which spots are currently occupied

System issues a ticket at entry:
- Active tickets to validate on exit

Calculates fee based on time spent (hourly):
- The hourly rate for pricing

who should track whether a spot is occupied? The spot itself or the parking lot?

One option is to add an isOccupied boolean field to ParkingSpot. You mark it true when someone parks there, and false when they leave.
The search logic is simple and direct. The spot “knows” its own state, you just ask it. The code also reads cleanly which is nice, “find a spot that’s free and matches the type.”
Challenge:
The main drawback of this approach is data duplication. Occupancy is represented in two different places: a boolean flag on the parking spot and the activeTickets map (since a ticket pointing to a spot implies that the spot is taken). This means you now have two sources of truth, and it becomes your responsibility to keep them in sync. If you miss an update in any code path, the state can drift, leading to issues such as allocating the same spot more than once.
There is also a more fundamental design concern around what “occupied” actually means. Attributes like a spot’s ID and size are intrinsic and never change — they define the spot itself. Occupancy, on the other hand, is temporary and contextual. It does not describe the spot; it describes the current state of the parking lot.
For example, the moment a ticket is issued at the entry gate, the spot should be considered occupied, even though the car has not yet parked. Similarly, when a car leaves its spot, the spot is still considered occupied until the ticket is closed at the exit. This makes it clear that occupancy is tied to ticket assignment, not to the physical parking spot.
You could skip storing occupancy entirely. A spot is occupied if and only if an active ticket references it. Compute this on demand:
For a small parking lot (~200 spots, ~100 active tickets), this works fine. The computation is cheap, it’s a few microseconds per entry.
But you’re computing the same thing repeatedly. Every time someone enters, you scan all tickets to build the occupied set. And if you later add features that need concurrency (multiple entrances processing cars simultaneously), you’d need to lock the entire activeTickets map during both reads and writes. You can’t have one thread iterating over tickets while another modifies the map.
This isn’t a problem for the base requirements, but it creates headaches if you extend the system later. Despite these considerations, this is conceptually the cleanest approach — there’s no denormalization at all. For small to medium systems, the simplicity might outweigh the concurrency and performance considerations.
To build on the previous approach, ParkingSpot can still stay a pure data holder and ParkingLot can track which spots are occupied using a separate Set.
Yes, the Set is technically redundant with the ticket data. You could derive occupancy from tickets. But the Set serves as a maintained index. Like a database index, it must stay in sync with the source data, but the trade-off is worth it.
We’re choosing the indexed approach for this problem. It keeps occupancy tracking centralized, gives us O(1) lookups, and sets us up well for the concurrency patterns we’ll discuss in the extensibility section. But as we discussed above, the other approaches are equally valid for different priorities.

ParkingLot
- spots: List<ParkingSpot>
  -occupiedSpotIds: Set<String>
  -activeTickets: Map<String, Ticket>
  -hourlyRateCents: long
+ ParkingLot(spots, hourlyRateCents)
+ enter(vehicleType) -> Ticket
+ exit(ticketId) -> long

ParkingSpot
ParkingSpot:
— id: String
— spotType: SpotType

+ ParkingSpot(id, spotType)
+ getSpotType() -> SpotType
+ getId() -> String

enum SpotType:
MOTORCYCLE
CAR
LARGE

enum VehicleType:
MOTORCYCLE
CAR
LARGE

Ticket
What Ticket must track

- When a vehicle exits, user provides ticket ID
  Ticket ID string

- Frees the spot for next use
  Which spot the vehicle is in

- System supports three vehicle types
  Type of vehicle (not used in base pricing, but stored for per-type pricing extension)

- Calculates fee based on time spent
  When they entered (needed for fee calculation)

Ticket
- id: string
- entryTime :timestamp
- exitTime: timestamp
- VehicalType : VehicleType
- status : string
- spotId : string
+ Ticket(id, spotId, vehicleType, entryTime)
+ getId() -> String
+ getSpotId() -> String
+ getVehicleType() -> VehicleType
+ getEntryTime() -> long2