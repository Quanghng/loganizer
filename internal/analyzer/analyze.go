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
    LogID:        res.InputTarget.Id,
    FilePath:     res.InputTarget.Path,
    Status:       res.status, // assuming Status is exported
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