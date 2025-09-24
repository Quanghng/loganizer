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
    LogID        string `json:"log_id"`
    FilePath     string `json:"file_path"`
    Status       string `json:"status"`
    Message      string `json:"message"`
    ErrorDetails string `json:"error_details,omitempty"`
}

type AnalyzeResult struct {
	InputTarget config.InputTarget
	status			string
	Err				error
}

func AnalyzeLogSync (target config.InputTarget) AnalyzeResult {
	ms := rand.Intn(151) + 50                   
	d := time.Duration(ms) * time.Millisecond
	time.Sleep(d)

	// Verifier si le fichier existe
	if _, err := os.Stat(target.Path); err != nil {
		if os.IsNotExist(err) {
			return AnalyzeResult{
				InputTarget: target,
				status:      string(StatusFailed),
				Err:         &InaccessibleFileError{Err: err},
			}
		}
		// Autre erreur lors de l'accès au fichier
		return AnalyzeResult{
			InputTarget: target,
			status:      string(StatusFailed),
			Err:         &InaccessibleFileError{Err: err},
		}
	}
	
	// Success
	return AnalyzeResult{
		InputTarget: target,
		status:      string(StatusOK),
		Err:         nil,
	}
}

func ConvertToReportEntry(res AnalyzeResult) ReportEntry {
	report := ReportEntry{
    LogID:        res.InputTarget.Id,
    FilePath:     res.InputTarget.Path,
    Status:       res.status,
    Message:      "Analyse terminée avec succès.",
	}


	if res.Err != nil {
		var inaccessible *InaccessibleFileError
		var parsing *ParsingError
		if errors.As(res.Err, &inaccessible) {
			report.Status = "FAILED"
			report.Message = "Fichier introuvable."
			report.ErrorDetails = fmt.Sprintf("Inaccessible: %v", inaccessible.Err)
		} else if errors.As(res.Err, &parsing) {
			report.Status = "FAILED"
			report.Message = "Parsing error."
			report.ErrorDetails = fmt.Sprintf("Parsing error: %v", parsing.Err)
		} else {
			report.Status = "FAILED"
			report.Message = "Generic error."
			report.ErrorDetails = fmt.Sprintf("Generic error: %v", res.Err)
		}
	}
	
	return report
}