package application

import (
	"bytes"
	"fmt"
	"strings"
	"time"
	"wedding-web/internal/domain"
)

// CalendarService génère des fichiers ICS
type CalendarService struct{}

// NewCalendarService crée un nouveau service calendar
func NewCalendarService() *CalendarService {
	return &CalendarService{}
}

// GenerateICS génère un fichier .ics pour le planning
func (s *CalendarService) GenerateICS(planning *domain.Planning) ([]byte, error) {
	var buf bytes.Buffer

	// En-tête ICS
	buf.WriteString("BEGIN:VCALENDAR\r\n")
	buf.WriteString("VERSION:2.0\r\n")
	buf.WriteString("PRODID:-//Wedding Web//FR\r\n")
	buf.WriteString("CALSCALE:GREGORIAN\r\n")
	buf.WriteString("METHOD:PUBLISH\r\n")
	buf.WriteString("X-WR-CALNAME:Mariage\r\n")
	buf.WriteString("X-WR-TIMEZONE:Europe/Paris\r\n")

	// Ajouter chaque événement
	for _, event := range planning.Events {
		buf.WriteString("BEGIN:VEVENT\r\n")
		buf.WriteString(fmt.Sprintf("UID:%s@wedding-web\r\n", generateEventUID(event)))
		buf.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", formatICSDate(time.Now())))
		buf.WriteString(fmt.Sprintf("DTSTART:%s\r\n", formatICSDate(event.StartTime)))
		buf.WriteString(fmt.Sprintf("DTEND:%s\r\n", formatICSDate(event.EndTime)))
		buf.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", escapeICS(event.Title)))
		buf.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", escapeICS(event.Description)))
		buf.WriteString(fmt.Sprintf("LOCATION:%s\\, %s\r\n", escapeICS(event.Location), escapeICS(event.Address)))
		buf.WriteString("STATUS:CONFIRMED\r\n")
		buf.WriteString("SEQUENCE:0\r\n")
		buf.WriteString("END:VEVENT\r\n")
	}

	buf.WriteString("END:VCALENDAR\r\n")

	return buf.Bytes(), nil
}

// formatICSDate formate une date pour ICS (format: 20260711T143000Z)
func formatICSDate(t time.Time) string {
	return t.UTC().Format("20060102T150405Z")
}

// escapeICS échappe les caractères spéciaux pour ICS
func escapeICS(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, ",", "\\,")
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}

// generateEventUID génère un UID unique pour un événement
func generateEventUID(event domain.PlanningEvent) string {
	return fmt.Sprintf("%s-%d",
		strings.ReplaceAll(strings.ToLower(event.Title), " ", "-"),
		event.StartTime.Unix())
}
