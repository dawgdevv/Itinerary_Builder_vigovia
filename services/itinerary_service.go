package services

import (
	"strings"
	"time"

	"vigovia-task/models"
	"vigovia-task/storage"
	"vigovia-task/utils"
)

// ItineraryService handles business logic for itineraries
type ItineraryService struct {
	store *storage.MemoryStore
}

// NewItineraryService creates a new instance of ItineraryService
func NewItineraryService(store *storage.MemoryStore) *ItineraryService {
	return &ItineraryService{
		store: store,
	}
}

// CreateItinerary creates a new itinerary
func (is *ItineraryService) CreateItinerary(req *models.CreateItineraryRequest) (*models.Itinerary, error) {
	// Validate the request
	if err := utils.ValidateItinerary(req); err != nil {
		return nil, err
	}

	now := time.Now()
	itinerary := &models.Itinerary{
		ID:          utils.GenerateID(),
		UserID:      strings.TrimSpace(req.UserID),
		Title:       req.Title,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Location:    req.Location,
		Hotels:      req.Hotels,
		Flights:     req.Flights,
		Transfers:   req.Transfers,
		Days:        req.Days,
		PaymentPlan: req.PaymentPlan,
		Inclusions:  req.Inclusions,
		Exclusions:  req.Exclusions,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Normalise activity periods immediately for consistent downstream usage.
	for i := range itinerary.Days {
		for j := range itinerary.Days[i].Activities {
			itinerary.Days[i].Activities[j] = normalizeActivity(itinerary.Days[i].Activities[j])
		}
	}

	if err := is.store.Create(itinerary); err != nil {
		return nil, err
	}

	return itinerary, nil
}

// GetItinerary retrieves an itinerary by ID
func (is *ItineraryService) GetItinerary(id string) (*models.Itinerary, error) {
	return is.store.GetByID(id)
}

// ListItineraries retrieves all itineraries
func (is *ItineraryService) ListItineraries() []*models.Itinerary {
	return is.store.GetAll()
}

// UpdateItinerary updates an existing itinerary
func (is *ItineraryService) UpdateItinerary(id string, req *models.UpdateItineraryRequest) (*models.Itinerary, error) {
	// Get the existing itinerary
	itinerary, err := is.store.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if strings.TrimSpace(req.UserID) != "" {
		itinerary.UserID = strings.TrimSpace(req.UserID)
	}
	if req.Title != "" {
		itinerary.Title = req.Title
	}
	if req.Description != "" {
		itinerary.Description = req.Description
	}
	if !req.StartDate.IsZero() {
		itinerary.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		itinerary.EndDate = req.EndDate
	}
	if req.Location != "" {
		itinerary.Location = req.Location
	}
	if req.Hotels != nil {
		if len(req.Hotels) == 0 {
			return nil, utils.NewValidationError("at least one hotel is required")
		}
		for _, hotel := range req.Hotels {
			if err := utils.ValidateHotel(&hotel); err != nil {
				return nil, err
			}
		}
		itinerary.Hotels = req.Hotels
	}
	if req.Flights != nil {
		if len(req.Flights) == 0 {
			return nil, utils.NewValidationError("at least one flight is required")
		}
		for _, flight := range req.Flights {
			if err := utils.ValidateFlight(&flight); err != nil {
				return nil, err
			}
		}
		itinerary.Flights = req.Flights
	}
	if req.Transfers != nil {
		if len(req.Transfers) == 0 {
			return nil, utils.NewValidationError("at least one transfer is required")
		}
		for _, transfer := range req.Transfers {
			if err := utils.ValidateTransfer(&transfer); err != nil {
				return nil, err
			}
		}
		itinerary.Transfers = req.Transfers
	}
	if req.Days != nil {
		if len(req.Days) == 0 {
			return nil, utils.NewValidationError("at least one day plan is required")
		}
		for _, day := range req.Days {
			if err := utils.ValidateDayPlan(&day); err != nil {
				return nil, err
			}
		}
		itinerary.Days = req.Days
		for i := range itinerary.Days {
			for j := range itinerary.Days[i].Activities {
				itinerary.Days[i].Activities[j] = normalizeActivity(itinerary.Days[i].Activities[j])
			}
		}
	}
	if req.PaymentPlan != nil {
		if len(req.PaymentPlan) == 0 {
			return nil, utils.NewValidationError("at least one payment installment is required")
		}
		for _, installment := range req.PaymentPlan {
			if err := utils.ValidatePaymentInstallment(&installment); err != nil {
				return nil, err
			}
		}
		itinerary.PaymentPlan = req.PaymentPlan
	}
	if req.Inclusions != nil {
		if err := utils.ValidateStringList(req.Inclusions, "inclusion"); err != nil {
			return nil, err
		}
		itinerary.Inclusions = req.Inclusions
	}
	if req.Exclusions != nil {
		if err := utils.ValidateStringList(req.Exclusions, "exclusion"); err != nil {
			return nil, err
		}
		itinerary.Exclusions = req.Exclusions
	}

	itinerary.UpdatedAt = time.Now()

	// Update in storage
	if err := is.store.Update(id, itinerary); err != nil {
		return nil, err
	}

	return itinerary, nil
}

// DeleteItinerary deletes an itinerary
func (is *ItineraryService) DeleteItinerary(id string) error {
	return is.store.Delete(id)
}

// AddActivity adds an activity to a specific day
func (is *ItineraryService) AddActivity(itineraryID string, dayNumber int, activity *models.Activity) (*models.Itinerary, error) {
	if err := utils.ValidateActivity(activity); err != nil {
		return nil, err
	}

	itinerary, err := is.store.GetByID(itineraryID)
	if err != nil {
		return nil, err
	}

	// Find the day and add the activity
	for i := range itinerary.Days {
		if itinerary.Days[i].DayNumber == dayNumber {
			itinerary.Days[i].Activities = append(itinerary.Days[i].Activities, normalizeActivity(*activity))
			itinerary.UpdatedAt = time.Now()
			if err := is.store.Update(itineraryID, itinerary); err != nil {
				return nil, err
			}
			return itinerary, nil
		}
	}

	return nil, utils.NewValidationError("day not found in itinerary")
}

func normalizeActivity(activity models.Activity) models.Activity {
	activity.Period = strings.ToLower(activity.Period)
	return activity
}
