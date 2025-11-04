package services

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"vigovia-task/models"

	"github.com/jung-kurt/gofpdf"
)

// PDFService handles PDF generation for itineraries
type PDFService struct{}

var (
	spacePattern   = regexp.MustCompile(`\s+`)
	invalidPattern = regexp.MustCompile(`[^a-z0-9-_]+`)
	periodOrder    = []string{models.ActivityPeriodMorning, models.ActivityPeriodAfternoon, models.ActivityPeriodEvening}
	periodDisplay  = map[string]string{
		models.ActivityPeriodMorning:   "Morning",
		models.ActivityPeriodAfternoon: "Afternoon",
		models.ActivityPeriodEvening:   "Evening",
	}
)

// NewPDFService creates a new instance of PDFService
func NewPDFService() *PDFService {
	return &PDFService{}
}

// GeneratePDF generates a professional PDF document for an itinerary
func (ps *PDFService) GeneratePDF(itinerary *models.Itinerary) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetLeftMargin(15)
	pdf.SetRightMargin(15)

	// Header with background color
	ps.addHeader(pdf, itinerary)

	// Metadata section
	ps.addMetadataSection(pdf, itinerary)

	if len(itinerary.Hotels) > 0 {
		ps.addHotelsSection(pdf, itinerary.Hotels)
	}

	if len(itinerary.Flights) > 0 {
		ps.addFlightsSection(pdf, itinerary.Flights)
	}

	if len(itinerary.Transfers) > 0 {
		ps.addTransfersSection(pdf, itinerary.Transfers)
	}

	// Description section
	if itinerary.Description != "" {
		ps.addDescriptionSection(pdf, itinerary.Description)
	}

	// Itinerary details section
	if len(itinerary.Days) > 0 {
		ps.addItineraryDetailsSection(pdf, itinerary.Days)
	} else {
		// If no days, add empty state message
		pdf.SetFont("Arial", "I", 11)
		pdf.SetTextColor(100, 100, 100)
		pdf.Cell(0, 10, "No days planned yet")
	}

	if len(itinerary.PaymentPlan) > 0 {
		ps.addPaymentPlanSection(pdf, itinerary.PaymentPlan)
	}

	if len(itinerary.Inclusions) > 0 || len(itinerary.Exclusions) > 0 {
		ps.addInclusionsExclusionsSection(pdf, itinerary.Inclusions, itinerary.Exclusions)
	}

	// Footer
	ps.addFooter(pdf)

	// Generate PDF bytes
	buf := new(bytes.Buffer)
	if err := pdf.Output(buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	pdfBytes := buf.Bytes()
	if err := ps.savePDFToOutput(pdfBytes, itinerary); err != nil {
		return nil, fmt.Errorf("failed to save PDF: %w", err)
	}

	return pdfBytes, nil
}

// addHeader adds a professional header to the PDF
func (ps *PDFService) addHeader(pdf *gofpdf.Fpdf, itinerary *models.Itinerary) {
	// Minimal header layout keeps document clean
	pdf.SetXY(15, 18)
	pdf.SetFont("Arial", "B", 19)
	pdf.SetTextColor(34, 34, 34)
	pdf.CellFormat(0, 9, "Itinerary Plan", "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(90, 90, 90)
	pdf.CellFormat(0, 6, itinerary.Title, "", 1, "L", false, 0, "")
	pdf.Ln(2)

	ps.drawDivider(pdf)
	pdf.Ln(8)
}

// addMetadataSection adds trip information
func (ps *PDFService) addMetadataSection(pdf *gofpdf.Fpdf, itinerary *models.Itinerary) {
	ps.addSectionHeader(pdf, "Trip Information")

	ps.addLabelValue(pdf, "User ID", itinerary.UserID)
	ps.addLabelValue(pdf, "Location", itinerary.Location)
	ps.addLabelValue(pdf, "Start Date", formatDate(itinerary.StartDate))
	ps.addLabelValue(pdf, "End Date", formatDate(itinerary.EndDate))
	ps.addLabelValue(pdf, "Duration", fmt.Sprintf("%d Days", len(itinerary.Days)))
	pdf.Ln(6)
}

// addDescriptionSection adds the itinerary description
func (ps *PDFService) addDescriptionSection(pdf *gofpdf.Fpdf, description string) {
	ps.addSectionHeader(pdf, "Overview")

	pdf.SetTextColor(50, 50, 50)
	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(0, 6, description, "", "L", false)
	pdf.Ln(8)
}

// addItineraryDetailsSection adds detailed day-by-day itinerary
func (ps *PDFService) addItineraryDetailsSection(pdf *gofpdf.Fpdf, days []models.DayPlan) {
	ps.addSectionHeader(pdf, "Day-by-Day Itinerary")

	for _, day := range days {
		ps.addDaySection(pdf, day)
	}
}

// addDaySection adds a single day's information
func (ps *PDFService) addDaySection(pdf *gofpdf.Fpdf, day models.DayPlan) {
	// Day header
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(41, 128, 185)
	pdf.SetX(15)
	pdf.CellFormat(0, 7, fmt.Sprintf("Day %d: %s", day.DayNumber, day.Title), "", 1, "L", false, 0, "")

	// Date line
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(110, 110, 110)
	pdf.SetX(15)
	pdf.CellFormat(0, 5, fmt.Sprintf("Date: %s", formatDate(day.Date)), "", 1, "L", false, 0, "")
	pdf.Ln(1)

	// Activities
	if len(day.Activities) > 0 {
		grouped := make(map[string][]models.Activity)
		for _, activity := range day.Activities {
			periodKey := strings.ToLower(activity.Period)
			grouped[periodKey] = append(grouped[periodKey], activity)
		}

		for _, period := range periodOrder {
			activities := grouped[period]
			if len(activities) == 0 {
				continue
			}

			label := periodDisplay[period]
			if label == "" {
				label = toTitleCase(period)
			}

			pdf.SetFont("Arial", "B", 11)
			pdf.SetTextColor(70, 70, 70)
			pdf.SetX(18)
			pdf.CellFormat(0, 6, label+" Session", "", 1, "L", false, 0, "")
			pdf.Ln(1)

			for idx, activity := range activities {
				pdf.SetFont("Arial", "B", 10)
				pdf.SetTextColor(52, 152, 219)
				pdf.SetX(20)
				pdf.CellFormat(28, 5, activity.Time, "", 0, "L", false, 0, "")

				pdf.SetFont("Arial", "B", 10)
				pdf.SetTextColor(40, 40, 40)
				pdf.CellFormat(0, 5, activity.Title, "", 1, "L", false, 0, "")

				pdf.SetFont("Arial", "", 10)
				pdf.SetTextColor(90, 90, 90)
				if activity.Description != "" {
					pdf.SetX(25)
					pdf.MultiCell(0, 5, "Description: "+activity.Description, "", "L", false)
				}
				if activity.Location != "" {
					pdf.SetX(25)
					pdf.MultiCell(0, 5, "Location: "+activity.Location, "", "L", false)
				}
				if activity.Duration != "" {
					pdf.SetX(25)
					pdf.MultiCell(0, 5, "Duration: "+activity.Duration, "", "L", false)
				}

				if idx < len(activities)-1 {
					pdf.Ln(2)
				}
			}

			pdf.Ln(3)
		}

		for period, activities := range grouped {
			skip := false
			for _, ordered := range periodOrder {
				if period == ordered {
					skip = true
					break
				}
			}
			if skip || len(activities) == 0 {
				continue
			}

			label := toTitleCase(period)
			pdf.SetFont("Arial", "B", 11)
			pdf.SetTextColor(70, 70, 70)
			pdf.SetX(18)
			pdf.CellFormat(0, 6, label+" Session", "", 1, "L", false, 0, "")
			pdf.Ln(1)

			for idx, activity := range activities {
				pdf.SetFont("Arial", "B", 10)
				pdf.SetTextColor(52, 152, 219)
				pdf.SetX(20)
				pdf.CellFormat(28, 5, activity.Time, "", 0, "L", false, 0, "")

				pdf.SetFont("Arial", "B", 10)
				pdf.SetTextColor(40, 40, 40)
				pdf.CellFormat(0, 5, activity.Title, "", 1, "L", false, 0, "")

				pdf.SetFont("Arial", "", 10)
				pdf.SetTextColor(90, 90, 90)
				if activity.Description != "" {
					pdf.SetX(25)
					pdf.MultiCell(0, 5, "Description: "+activity.Description, "", "L", false)
				}
				if activity.Location != "" {
					pdf.SetX(25)
					pdf.MultiCell(0, 5, "Location: "+activity.Location, "", "L", false)
				}
				if activity.Duration != "" {
					pdf.SetX(25)
					pdf.MultiCell(0, 5, "Duration: "+activity.Duration, "", "L", false)
				}

				if idx < len(activities)-1 {
					pdf.Ln(2)
				}
			}

			pdf.Ln(3)
		}
	} else {
		pdf.SetFont("Arial", "I", 10)
		pdf.SetTextColor(150, 150, 150)
		pdf.SetX(25)
		pdf.Cell(0, 6, "No activities planned for this day")
		pdf.Ln(8)
	}

	// Space between days
	pdf.Ln(3)
	pdf.SetDrawColor(220, 220, 220)
	pdf.Line(15, pdf.GetY(), 195, pdf.GetY())
	pdf.Ln(8)
}

func (ps *PDFService) addSectionHeader(pdf *gofpdf.Fpdf, title string) {
	if title == "" {
		return
	}

	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(41, 128, 185)
	pdf.SetX(15)
	pdf.CellFormat(0, 7, title, "", 1, "L", false, 0, "")
	ps.drawDivider(pdf)
	pdf.Ln(6)
	pdf.SetTextColor(0, 0, 0)
}

func (ps *PDFService) addLabelValue(pdf *gofpdf.Fpdf, label, value string) {
	if value == "" {
		return
	}

	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(70, 70, 70)
	pdf.SetX(15)
	pdf.CellFormat(30, 5, fmt.Sprintf("%s:", label), "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(35, 35, 35)
	pdf.MultiCell(0, 5, value, "", "L", false)
	pdf.Ln(1)
}

func (ps *PDFService) drawDivider(pdf *gofpdf.Fpdf) {
	y := pdf.GetY()
	pdf.SetDrawColor(210, 210, 210)
	pdf.Line(15, y, 195, y)
}

func (ps *PDFService) addHotelsSection(pdf *gofpdf.Fpdf, hotels []models.Hotel) {
	if len(hotels) == 0 {
		return
	}

	ps.addSectionHeader(pdf, "Hotel Accommodations")

	for idx, hotel := range hotels {
		// Hotel name header with number
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(40, 40, 40)
		pdf.SetX(15)
		pdf.CellFormat(0, 6, fmt.Sprintf("Hotel %d: %s", idx+1, hotel.Name), "", 1, "L", false, 0, "")

		// Hotel details with labels
		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(90, 90, 90)
		
		if hotel.City != "" {
			pdf.SetX(20)
			pdf.SetFont("Arial", "B", 10)
			pdf.SetTextColor(70, 70, 70)
			pdf.CellFormat(25, 5, "Location:", "", 0, "L", false, 0, "")
			pdf.SetFont("Arial", "", 10)
			pdf.SetTextColor(90, 90, 90)
			pdf.CellFormat(0, 5, hotel.City, "", 1, "L", false, 0, "")
		}
		
		if !hotel.CheckIn.IsZero() {
			pdf.SetX(20)
			pdf.SetFont("Arial", "B", 10)
			pdf.SetTextColor(70, 70, 70)
			pdf.CellFormat(25, 5, "Check-in:", "", 0, "L", false, 0, "")
			pdf.SetFont("Arial", "", 10)
			pdf.SetTextColor(90, 90, 90)
			pdf.CellFormat(0, 5, formatDate(hotel.CheckIn), "", 1, "L", false, 0, "")
		}
		
		if !hotel.CheckOut.IsZero() {
			pdf.SetX(20)
			pdf.SetFont("Arial", "B", 10)
			pdf.SetTextColor(70, 70, 70)
			pdf.CellFormat(25, 5, "Check-out:", "", 0, "L", false, 0, "")
			pdf.SetFont("Arial", "", 10)
			pdf.SetTextColor(90, 90, 90)
			pdf.CellFormat(0, 5, formatDate(hotel.CheckOut), "", 1, "L", false, 0, "")
		}
		
		if hotel.Nights > 0 {
			pdf.SetX(20)
			pdf.SetFont("Arial", "B", 10)
			pdf.SetTextColor(70, 70, 70)
			pdf.CellFormat(25, 5, "Duration:", "", 0, "L", false, 0, "")
			pdf.SetFont("Arial", "", 10)
			pdf.SetTextColor(90, 90, 90)
			nightsText := fmt.Sprintf("%d night", hotel.Nights)
			if hotel.Nights > 1 {
				nightsText += "s"
			}
			pdf.CellFormat(0, 5, nightsText, "", 1, "L", false, 0, "")
		}
		
		if idx < len(hotels)-1 {
			pdf.Ln(3)
		}
	}

	pdf.Ln(4)
}

func (ps *PDFService) addFlightsSection(pdf *gofpdf.Fpdf, flights []models.Flight) {
	if len(flights) == 0 {
		return
	}

	ps.addSectionHeader(pdf, "Flight Details")

	for idx, flight := range flights {
		// Flight header with number
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(40, 40, 40)
		pdf.SetX(15)
		
		flightHeader := fmt.Sprintf("Flight %d", idx+1)
		if flight.Airline != "" && flight.FlightNumber != "" {
			flightHeader = fmt.Sprintf("Flight %d: %s %s", idx+1, flight.Airline, flight.FlightNumber)
		} else if flight.Airline != "" {
			flightHeader = fmt.Sprintf("Flight %d: %s", idx+1, flight.Airline)
		} else if flight.FlightNumber != "" {
			flightHeader = fmt.Sprintf("Flight %d: %s", idx+1, flight.FlightNumber)
		}
		
		pdf.CellFormat(0, 6, flightHeader, "", 1, "L", false, 0, "")

		// Departure details
		if flight.DepartureCity != "" || flight.DepartureAirport != "" || !flight.DepartureTime.IsZero() {
			pdf.SetFont("Arial", "B", 10)
			pdf.SetTextColor(70, 70, 70)
			pdf.SetX(20)
			pdf.CellFormat(0, 5, "Departure", "", 1, "L", false, 0, "")
			
			if flight.DepartureCity != "" {
				pdf.SetX(25)
				pdf.SetFont("Arial", "B", 10)
				pdf.SetTextColor(70, 70, 70)
				pdf.CellFormat(20, 5, "City:", "", 0, "L", false, 0, "")
				pdf.SetFont("Arial", "", 10)
				pdf.SetTextColor(90, 90, 90)
				pdf.CellFormat(0, 5, flight.DepartureCity, "", 1, "L", false, 0, "")
			}
			
			if flight.DepartureAirport != "" {
				pdf.SetX(25)
				pdf.SetFont("Arial", "B", 10)
				pdf.SetTextColor(70, 70, 70)
				pdf.CellFormat(20, 5, "Airport:", "", 0, "L", false, 0, "")
				pdf.SetFont("Arial", "", 10)
				pdf.SetTextColor(90, 90, 90)
				pdf.CellFormat(0, 5, flight.DepartureAirport, "", 1, "L", false, 0, "")
			}
			
			if !flight.DepartureTime.IsZero() {
				pdf.SetX(25)
				pdf.SetFont("Arial", "B", 10)
				pdf.SetTextColor(70, 70, 70)
				pdf.CellFormat(20, 5, "Time:", "", 0, "L", false, 0, "")
				pdf.SetFont("Arial", "", 10)
				pdf.SetTextColor(90, 90, 90)
				pdf.CellFormat(0, 5, formatDateTime(flight.DepartureTime), "", 1, "L", false, 0, "")
			}
		}

		pdf.Ln(1)

		// Arrival details
		if flight.ArrivalCity != "" || flight.ArrivalAirport != "" || !flight.ArrivalTime.IsZero() {
			pdf.SetFont("Arial", "B", 10)
			pdf.SetTextColor(70, 70, 70)
			pdf.SetX(20)
			pdf.CellFormat(0, 5, "Arrival", "", 1, "L", false, 0, "")
			
			if flight.ArrivalCity != "" {
				pdf.SetX(25)
				pdf.SetFont("Arial", "B", 10)
				pdf.SetTextColor(70, 70, 70)
				pdf.CellFormat(20, 5, "City:", "", 0, "L", false, 0, "")
				pdf.SetFont("Arial", "", 10)
				pdf.SetTextColor(90, 90, 90)
				pdf.CellFormat(0, 5, flight.ArrivalCity, "", 1, "L", false, 0, "")
			}
			
			if flight.ArrivalAirport != "" {
				pdf.SetX(25)
				pdf.SetFont("Arial", "B", 10)
				pdf.SetTextColor(70, 70, 70)
				pdf.CellFormat(20, 5, "Airport:", "", 0, "L", false, 0, "")
				pdf.SetFont("Arial", "", 10)
				pdf.SetTextColor(90, 90, 90)
				pdf.CellFormat(0, 5, flight.ArrivalAirport, "", 1, "L", false, 0, "")
			}
			
			if !flight.ArrivalTime.IsZero() {
				pdf.SetX(25)
				pdf.SetFont("Arial", "B", 10)
				pdf.SetTextColor(70, 70, 70)
				pdf.CellFormat(20, 5, "Time:", "", 0, "L", false, 0, "")
				pdf.SetFont("Arial", "", 10)
				pdf.SetTextColor(90, 90, 90)
				pdf.CellFormat(0, 5, formatDateTime(flight.ArrivalTime), "", 1, "L", false, 0, "")
			}
		}
		
		if idx < len(flights)-1 {
			pdf.Ln(4)
		}
	}

	pdf.Ln(4)
}

func (ps *PDFService) addTransfersSection(pdf *gofpdf.Fpdf, transfers []models.Transfer) {
	if len(transfers) == 0 {
		return
	}

	ps.addSectionHeader(pdf, "Transfers")

	for _, transfer := range transfers {
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(40, 40, 40)
		pdf.SetX(15)
		pdf.CellFormat(0, 6, toTitleCase(transfer.Mode)+" Transfer", "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(90, 90, 90)
		pdf.SetX(15)
		pdf.CellFormat(0, 5, fmt.Sprintf("Pickup: %s   Drop-off: %s   Time: %s", transfer.Pickup, transfer.Dropoff, transfer.PickupTime), "", 1, "L", false, 0, "")
		if strings.TrimSpace(transfer.Notes) != "" {
			pdf.SetX(15)
			pdf.MultiCell(0, 5, "Notes: "+transfer.Notes, "", "L", false)
		}
		pdf.Ln(2)
	}

	pdf.Ln(4)
}

func (ps *PDFService) addPaymentPlanSection(pdf *gofpdf.Fpdf, plan []models.PaymentInstallment) {
	if len(plan) == 0 {
		return
	}

	ps.addSectionHeader(pdf, "Payment Plan")

	ordered := make([]models.PaymentInstallment, len(plan))
	copy(ordered, plan)
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].InstallmentNumber < ordered[j].InstallmentNumber
	})

	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(60, 60, 60)
	pdf.SetX(15)
	pdf.CellFormat(30, 6, "Installment", "", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, "Amount", "", 0, "L", false, 0, "")
	pdf.CellFormat(35, 6, "Due Date", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Status", "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(40, 40, 40)
	for _, installment := range ordered {
		status := strings.TrimSpace(installment.Status)
		if status == "" {
			status = "Pending"
		}

		pdf.SetX(15)
		pdf.CellFormat(30, 6, fmt.Sprintf("#%d", installment.InstallmentNumber), "", 0, "L", false, 0, "")
		pdf.CellFormat(35, 6, formatAmount(installment.Amount, installment.Currency), "", 0, "L", false, 0, "")
		pdf.CellFormat(35, 6, formatDate(installment.DueDate), "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 6, status, "", 1, "L", false, 0, "")
	}

	pdf.Ln(6)
}

func (ps *PDFService) addInclusionsExclusionsSection(pdf *gofpdf.Fpdf, inclusions, exclusions []string) {
	if len(inclusions) == 0 && len(exclusions) == 0 {
		return
	}

	ps.addSectionHeader(pdf, "Inclusions and Exclusions")

	if len(inclusions) > 0 {
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(70, 70, 70)
		pdf.SetX(15)
		pdf.CellFormat(0, 6, "Inclusions", "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(90, 90, 90)
		for _, item := range inclusions {
			pdf.SetX(20)
			pdf.CellFormat(0, 5, "- "+item, "", 1, "L", false, 0, "")
		}
		pdf.Ln(2)
	}

	if len(exclusions) > 0 {
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(70, 70, 70)
		pdf.SetX(15)
		pdf.CellFormat(0, 6, "Exclusions", "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "", 10)
		pdf.SetTextColor(90, 90, 90)
		for _, item := range exclusions {
			pdf.SetX(20)
			pdf.CellFormat(0, 5, "- "+item, "", 1, "L", false, 0, "")
		}
		pdf.Ln(2)
	}

	pdf.Ln(4)
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("Jan 2, 2006")
}

func formatDateTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("Jan 2, 2006 15:04")
}

func formatAmount(amount float64, currency string) string {
	if currency == "" {
		currency = "USD"
	}
	return strings.ToUpper(currency) + " " + fmt.Sprintf("%.2f", amount)
}

func toTitleCase(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	lower := strings.ToLower(trimmed)
	return strings.ToUpper(lower[:1]) + lower[1:]
}

func (ps *PDFService) savePDFToOutput(pdfBytes []byte, itinerary *models.Itinerary) error {
	if err := os.MkdirAll("output", 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	fileName := ps.buildFileName(itinerary)
	filePath := filepath.Join("output", fileName)

	if err := os.WriteFile(filePath, pdfBytes, 0644); err != nil {
		return fmt.Errorf("write PDF to %s: %w", filePath, err)
	}

	return nil
}

func (ps *PDFService) buildFileName(itinerary *models.Itinerary) string {
	base := strings.TrimSpace(itinerary.Title)
	if base == "" {
		base = "itinerary"
	}

	switch {
	case !itinerary.StartDate.IsZero():
		base = fmt.Sprintf("%s_%s", base, itinerary.StartDate.Format("2006-01-02"))
	case !itinerary.EndDate.IsZero():
		base = fmt.Sprintf("%s_%s", base, itinerary.EndDate.Format("2006-01-02"))
	default:
		base = fmt.Sprintf("%s_%s", base, time.Now().Format("2006-01-02"))
	}

	sanitized := sanitizeFileName(base)
	if sanitized == "" {
		sanitized = "itinerary"
	}

	return sanitized + ".pdf"
}

func sanitizeFileName(input string) string {
	lowered := strings.ToLower(strings.TrimSpace(input))
	normalized := spacePattern.ReplaceAllString(lowered, "-")
	cleaned := invalidPattern.ReplaceAllString(normalized, "")
	return strings.Trim(cleaned, "-_")
}

// addFooter adds a professional footer
func (ps *PDFService) addFooter(pdf *gofpdf.Fpdf) {
	// Footer area - position footer at bottom
	pdf.SetXY(15, 270) // Approximate position for A4
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(150, 150, 150)
	pdf.SetDrawColor(220, 220, 220)
	pdf.Line(15, 268, 195, 268)

	pdf.SetXY(15, 270)
	pdf.Cell(0, 10, fmt.Sprintf("Generated by Vigovia Itinerary Builder | Page %d", pdf.PageNo()))
}
