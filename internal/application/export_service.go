package application

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExportService gère l'export des données
type ExportService struct {
	rsvpService *RSVPService
}

// NewExportService crée un nouveau service d'export
func NewExportService(rsvpService *RSVPService) *ExportService {
	return &ExportService{
		rsvpService: rsvpService,
	}
}

// ExportRSVPsToExcel exporte tous les RSVP vers un fichier Excel
func (s *ExportService) ExportRSVPsToExcel() (*excelize.File, error) {
	// Récupérer tous les RSVP
	rsvps, err := s.rsvpService.ListRSVPs()
	if err != nil {
		return nil, err
	}

	// Créer un nouveau fichier Excel
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Confirmations RSVP"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	// Définir les en-têtes
	headers := []string{"Prénom", "Nom", "Statut", "Adultes", "Enfants", "Total", "Allergies/Régimes", "Message", "Date"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Style pour les en-têtes
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 12},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E8E8E8"}, Pattern: 1},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", "I1", headerStyle)
	}

	// Ajouter les données
	for i, rsvp := range rsvps {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), rsvp.FirstName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), rsvp.LastName)

		// Statut
		status := "Absent"
		if rsvp.WillAttend {
			status = "Présent"
		}
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), status)

		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), rsvp.AdultsCount)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), rsvp.ChildrenCount)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), rsvp.TotalGuests())
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), rsvp.Allergies)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), rsvp.Message)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), rsvp.SubmittedAt.Format("02/01/2006 15:04"))
	}

	// Ajouter une ligne de résumé
	summaryRow := len(rsvps) + 3
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", summaryRow), "TOTAL")

	// Formules pour les totaux
	f.SetCellFormula(sheetName, fmt.Sprintf("D%d", summaryRow), fmt.Sprintf("SUM(D2:D%d)", len(rsvps)+1))
	f.SetCellFormula(sheetName, fmt.Sprintf("E%d", summaryRow), fmt.Sprintf("SUM(E2:E%d)", len(rsvps)+1))
	f.SetCellFormula(sheetName, fmt.Sprintf("F%d", summaryRow), fmt.Sprintf("SUM(F2:F%d)", len(rsvps)+1))

	// Style pour la ligne de résumé
	summaryStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 11},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#D4E4F7"}, Pattern: 1},
	})
	if err == nil {
		f.SetCellStyle(sheetName, fmt.Sprintf("A%d", summaryRow), fmt.Sprintf("I%d", summaryRow), summaryStyle)
	}

	// Ajuster la largeur des colonnes
	f.SetColWidth(sheetName, "A", "B", 15)
	f.SetColWidth(sheetName, "C", "C", 12)
	f.SetColWidth(sheetName, "D", "F", 10)
	f.SetColWidth(sheetName, "G", "G", 30)
	f.SetColWidth(sheetName, "H", "H", 40)
	f.SetColWidth(sheetName, "I", "I", 18)

	// Activer les filtres
	f.AutoFilter(sheetName, "A1:I1", []excelize.AutoFilterOptions{})

	// Définir la feuille active
	f.SetActiveSheet(index)

	// Supprimer la feuille par défaut "Sheet1"
	f.DeleteSheet("Sheet1")

	return f, nil
}

// GetFileName génère un nom de fichier pour l'export
func (s *ExportService) GetFileName() string {
	return fmt.Sprintf("rsvp-mariage-%s.xlsx", time.Now().Format("2006-01-02"))
}
