package utils

import (
	"fmt"
	"strings"
	"time"

	"vigovia-task/models"
)

var validPeriods = map[string]struct{}{
	models.ActivityPeriodMorning:   {},
	models.ActivityPeriodAfternoon: {},
	models.ActivityPeriodEvening:   {},
}

// ValidateItinerary validates an itinerary
func ValidateItinerary(req *models.CreateItineraryRequest) error {
	if strings.TrimSpace(req.UserID) == "" {
		return NewValidationError("user_id is required")
	}

	if strings.TrimSpace(req.Title) == "" {
		return NewValidationError("title is required")
	}

	if strings.TrimSpace(req.Location) == "" {
		return NewValidationError("location is required")
	}

	if req.StartDate.IsZero() {
		return NewValidationError("start_date is required")
	}

	if req.EndDate.IsZero() {
		return NewValidationError("end_date is required")
	}

	if req.EndDate.Before(req.StartDate) {
		return NewValidationError("end_date must be after start_date")
	}

	if len(req.Hotels) == 0 {
		return NewValidationError("at least one hotel is required")
	}

	if len(req.Flights) == 0 {
		return NewValidationError("at least one flight is required")
	}

	if len(req.Transfers) == 0 {
		return NewValidationError("at least one transfer is required")
	}

	if len(req.Days) == 0 {
		return NewValidationError("at least one day plan is required")
	}

	if len(req.PaymentPlan) == 0 {
		return NewValidationError("at least one payment installment is required")
	}

	if err := ValidateStringList(req.Inclusions, "inclusion"); err != nil {
		return err
	}

	if err := ValidateStringList(req.Exclusions, "exclusion"); err != nil {
		return err
	}

	for _, hotel := range req.Hotels {
		if err := ValidateHotel(&hotel); err != nil {
			return err
		}
	}

	for _, flight := range req.Flights {
		if err := ValidateFlight(&flight); err != nil {
			return err
		}
	}

	for _, transfer := range req.Transfers {
		if err := ValidateTransfer(&transfer); err != nil {
			return err
		}
	}

	for _, day := range req.Days {
		if err := ValidateDayPlan(&day); err != nil {
			return err
		}
	}

	for _, installment := range req.PaymentPlan {
		if err := ValidatePaymentInstallment(&installment); err != nil {
			return err
		}
	}

	return nil
}

func ValidateStringList(values []string, label string) error {
	if len(values) == 0 {
		return NewValidationError(fmt.Sprintf("at least one %s is required", label))
	}

	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return NewValidationError(fmt.Sprintf("%s entries cannot be empty", label))
		}
	}

	return nil
}

// ValidateHotel ensures a hotel entry is well formed.
func ValidateHotel(hotel *models.Hotel) error {
	if strings.TrimSpace(hotel.Name) == "" {
		return NewValidationError("hotel name is required")
	}

	if strings.TrimSpace(hotel.City) == "" {
		return NewValidationError("hotel city is required")
	}

	if hotel.CheckIn.IsZero() {
		return NewValidationError("hotel check_in is required")
	}

	if hotel.CheckOut.IsZero() {
		return NewValidationError("hotel check_out is required")
	}

	if hotel.CheckOut.Before(hotel.CheckIn) {
		return NewValidationError("hotel check_out must be after check_in")
	}

	if hotel.Nights <= 0 {
		return NewValidationError("hotel nights must be greater than zero")
	}

	return nil
}

// ValidateFlight ensures a flight entry is well formed.
func ValidateFlight(flight *models.Flight) error {
	if strings.TrimSpace(flight.Airline) == "" {
		return NewValidationError("flight airline is required")
	}

	if strings.TrimSpace(flight.FlightNumber) == "" {
		return NewValidationError("flight number is required")
	}

	if strings.TrimSpace(flight.DepartureCity) == "" || strings.TrimSpace(flight.DepartureAirport) == "" {
		return NewValidationError("flight departure city and airport are required")
	}

	if flight.DepartureTime.IsZero() {
		return NewValidationError("flight departure_time is required")
	}

	if strings.TrimSpace(flight.ArrivalCity) == "" || strings.TrimSpace(flight.ArrivalAirport) == "" {
		return NewValidationError("flight arrival city and airport are required")
	}

	if flight.ArrivalTime.IsZero() {
		return NewValidationError("flight arrival_time is required")
	}

	if flight.ArrivalTime.Before(flight.DepartureTime) {
		return NewValidationError("flight arrival_time must be after departure_time")
	}

	return nil
}

// ValidateTransfer ensures transfer details contain the essentials.
func ValidateTransfer(transfer *models.Transfer) error {
	if strings.TrimSpace(transfer.Mode) == "" {
		return NewValidationError("transfer mode is required")
	}

	if strings.TrimSpace(transfer.Pickup) == "" {
		return NewValidationError("transfer pickup is required")
	}

	if strings.TrimSpace(transfer.Dropoff) == "" {
		return NewValidationError("transfer dropoff is required")
	}

	if strings.TrimSpace(transfer.PickupTime) == "" {
		return NewValidationError("transfer pickup_time is required")
	}

	return nil
}

// ValidateDayPlan ensures each day plan has mandatory details.
func ValidateDayPlan(day *models.DayPlan) error {
	if day.DayNumber <= 0 {
		return NewValidationError("day_number must be greater than zero")
	}

	if day.Date.IsZero() {
		return NewValidationError(fmt.Sprintf("date is required for day %d", day.DayNumber))
	}

	if strings.TrimSpace(day.Title) == "" {
		return NewValidationError(fmt.Sprintf("title is required for day %d", day.DayNumber))
	}

	if len(day.Activities) == 0 {
		return NewValidationError(fmt.Sprintf("at least one activity is required for day %d", day.DayNumber))
	}

	for _, activity := range day.Activities {
		if err := ValidateActivity(&activity); err != nil {
			return err
		}
	}

	return nil
}

// ValidateActivity validates an activity
func ValidateActivity(activity *models.Activity) error {
	if strings.TrimSpace(activity.Title) == "" {
		return NewValidationError("activity title is required")
	}

	if strings.TrimSpace(activity.Period) == "" {
		return NewValidationError("activity period is required")
	}

	if _, ok := validPeriods[strings.ToLower(activity.Period)]; !ok {
		return NewValidationError("activity period must be morning, afternoon, or evening")
	}

	if strings.TrimSpace(activity.Time) == "" {
		return NewValidationError("activity time is required")
	}

	if strings.TrimSpace(activity.Description) == "" {
		return NewValidationError("activity description is required")
	}

	if strings.TrimSpace(activity.Location) == "" {
		return NewValidationError("activity location is required")
	}

	return nil
}

// ValidatePaymentInstallment checks payment plan entries.
func ValidatePaymentInstallment(installment *models.PaymentInstallment) error {
	if installment.InstallmentNumber <= 0 {
		return NewValidationError("payment installment_number must be greater than zero")
	}

	if installment.Amount <= 0 {
		return NewValidationError("payment amount must be greater than zero")
	}

	if strings.TrimSpace(installment.Currency) == "" {
		return NewValidationError("payment currency is required")
	}

	if installment.DueDate.IsZero() {
		return NewValidationError("payment due_date is required")
	}

	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

// GenerateID generates a unique ID based on current timestamp
func GenerateID() string {
	return time.Now().Format("20060102150405")
}
