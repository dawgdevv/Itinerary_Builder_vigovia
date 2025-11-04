package models

import "time"

// ActivityPeriod constants ensure consistent categorisation of daily plans.
const (
	ActivityPeriodMorning   = "morning"
	ActivityPeriodAfternoon = "afternoon"
	ActivityPeriodEvening   = "evening"
)

// Itinerary represents a complete travel plan with all supporting sections.
type Itinerary struct {
	ID          string               `json:"id"`
	UserID      string               `json:"user_id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	StartDate   time.Time            `json:"start_date"`
	EndDate     time.Time            `json:"end_date"`
	Location    string               `json:"location"`
	Hotels      []Hotel              `json:"hotels"`
	Flights     []Flight             `json:"flights"`
	Transfers   []Transfer           `json:"transfers"`
	Days        []DayPlan            `json:"days"`
	PaymentPlan []PaymentInstallment `json:"payment_plan"`
	Inclusions  []string             `json:"inclusions"`
	Exclusions  []string             `json:"exclusions"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// Hotel captures accommodation details inside an itinerary.
type Hotel struct {
	Name     string    `json:"name"`
	City     string    `json:"city"`
	CheckIn  time.Time `json:"check_in"`
	CheckOut time.Time `json:"check_out"`
	Nights   int       `json:"nights"`
}

// Flight captures air travel segments of an itinerary.
type Flight struct {
	Airline          string    `json:"airline"`
	FlightNumber     string    `json:"flight_number"`
	DepartureCity    string    `json:"departure_city"`
	DepartureAirport string    `json:"departure_airport"`
	DepartureTime    time.Time `json:"departure_time"`
	ArrivalCity      string    `json:"arrival_city"`
	ArrivalAirport   string    `json:"arrival_airport"`
	ArrivalTime      time.Time `json:"arrival_time"`
}

// Transfer represents a ground transfer such as car or shuttle.
type Transfer struct {
	Mode       string `json:"mode"`
	Pickup     string `json:"pickup"`
	Dropoff    string `json:"dropoff"`
	PickupTime string `json:"pickup_time"`
	Notes      string `json:"notes"`
}

// PaymentInstallment describes a single entry in a payment plan.
type PaymentInstallment struct {
	InstallmentNumber int       `json:"installment_number"`
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency"`
	DueDate           time.Time `json:"due_date"`
	Status            string    `json:"status"`
}

// DayPlan represents a single day in the itinerary
type DayPlan struct {
	DayNumber  int        `json:"day_number"`
	Date       time.Time  `json:"date"`
	Title      string     `json:"title"`
	Activities []Activity `json:"activities"`
}

// Activity represents a single activity in a day plan
type Activity struct {
	Period      string `json:"period"`
	Time        string `json:"time"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Location    string `json:"location"`
	Duration    string `json:"duration"`
}

// CreateItineraryRequest is the request payload for creating an itinerary
type CreateItineraryRequest struct {
	UserID      string               `json:"user_id" binding:"required"`
	Title       string               `json:"title" binding:"required"`
	Description string               `json:"description"`
	StartDate   time.Time            `json:"start_date" binding:"required"`
	EndDate     time.Time            `json:"end_date" binding:"required"`
	Location    string               `json:"location" binding:"required"`
	Hotels      []Hotel              `json:"hotels" binding:"required"`
	Flights     []Flight             `json:"flights" binding:"required"`
	Transfers   []Transfer           `json:"transfers" binding:"required"`
	Days        []DayPlan            `json:"days" binding:"required"`
	PaymentPlan []PaymentInstallment `json:"payment_plan"`
	Inclusions  []string             `json:"inclusions"`
	Exclusions  []string             `json:"exclusions"`
}

// UpdateItineraryRequest is the request payload for updating an itinerary
type UpdateItineraryRequest struct {
	UserID      string               `json:"user_id"`
	Title       string               `json:"title"`
	Description string               `json:"description"`
	StartDate   time.Time            `json:"start_date"`
	EndDate     time.Time            `json:"end_date"`
	Location    string               `json:"location"`
	Hotels      []Hotel              `json:"hotels"`
	Flights     []Flight             `json:"flights"`
	Transfers   []Transfer           `json:"transfers"`
	Days        []DayPlan            `json:"days"`
	PaymentPlan []PaymentInstallment `json:"payment_plan"`
	Inclusions  []string             `json:"inclusions"`
	Exclusions  []string             `json:"exclusions"`
}
