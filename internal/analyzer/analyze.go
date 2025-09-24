package analyzer

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Quanghng/loganizer/internal/config"
)

// Statut possible
type Status string


const (
StatusOK Status = "OK"
StatusFailed Status = "FAILED"
)

type ReportEntry struct {
	log_id        string
	file_path     string
	status        string
	message       string // "OK", "Inaccessible", "Error"
	error_details string // Message d'erreur, omis si vide
}

type AnalyzeResult struct {
	InputTarget config.InputTarget
	status			string
	Err				error
}

func AnalyzeLogSync (target config.InputTarget) AnalyzeResult {
	ms := rand.Intn(151) + 50                     // random int between 50 and 200 inclusive
	d := time.Duration(ms) * time.Millisecond
	time.Sleep(d)

	if _, err := os.Stat(target.Path); err != nil {
		if os.IsNotExist(err) {
			return AnalyzeResult{}
		}
	}
	
	return AnalyzeResult{
		InputTarget: target,
		status: "OK",
	}
}

func ConvertToReportEntry(res AnalyzeResult) ReportEntry {
	report := ReportEntry{
		log_id:   res.InputTarget.Id,
		file_path:    res.InputTarget.Path,
		status:  res.status,
		message: "Analyse terminée avec succès.",
	}

	if res.Err != nil {
		var inaccessible *InaccessibleFileError
		var parsing *ParsingError
		if errors.As(res.Err, &inaccessible) {
			report.status = "FAILED"
			report.message = "Fichier introuvable."
			report.error_details = fmt.Sprintf("Inaccessible: %v", inaccessible.Err)
		} else if errors.As(res.Err, &parsing) {
			report.status = "FAILED"
			report.message = "Parsing error."
			report.error_details = fmt.Sprintf("Parsing error: %v", parsing.Err)
		} else {
			report.status = "FAILED"
			report.error_details = fmt.Sprintf("Erreur générique: %v", res.Err)
		}
	}
	
	return report
}