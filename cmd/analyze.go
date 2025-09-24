package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Quanghng/loganizer/internal/analyzer"
	"github.com/Quanghng/loganizer/internal/config"
	"github.com/Quanghng/loganizer/internal/reporter"
	"github.com/spf13/cobra"
)

var (
	cfgPath string
	outPath string
)


var analyzeCmd = &cobra.Command {
	Use: "analyze",
	Short: "Lancer l'analyse de logs à partir d'un fichier de configuration JSON",
	Run: func(cmd *cobra.Command, args []string) {

		if cfgPath == "" {
			fmt.Println("Erreur: le chemin du fichier d'entrée (--input) est obligatoire.")
			return
		}
		targets, err := config.LoadTargetsFromFile(cfgPath)
		if err != nil {
			fmt.Printf("Erreur lors du chargement des configs: %v\n", err)
			return
		}

		if len(targets) == 0 {
			fmt.Println("Aucune configs trouvée dans le fichier d'entrée.")
			return
		}

		var wg sync.WaitGroup
		resultsChan := make(chan analyzer.AnalyzeResult, len(targets))

		wg.Add(len(targets))
		for _, target := range targets {
			go func (t config.InputTarget) {
				defer wg.Done()
				result := analyzer.AnalyzeLogSync(t)
				resultsChan <- result
			}(target)
		}

		wg.Wait() 
		close(resultsChan)
		
		var finalReport []analyzer.ReportEntry
		for res := range resultsChan { // Récupérer tous les résultats du canal
			reportEntry := analyzer.ConvertToReportEntry(res)
			finalReport = append(finalReport, reportEntry)

			// Affichage immédiat comme avant
			if res.Err != nil {
				var inaccessible *analyzer.InaccessibleFileError
				var parsing *analyzer.ParsingError
				if  errors.As(res.Err, &inaccessible) {
					fmt.Printf("🚫 %s est inaccessible : %v\n", res.InputTarget.Id, inaccessible.Err)
				} else if errors.As(res.Err, &parsing) {
					fmt.Printf("erreur du parsing : %v\n", parsing.Err)
				} else {
					fmt.Printf("❌ %s (%s) : erreur - %v\n", res.InputTarget.Id, res.InputTarget.Path, res.Err)
				}
			} else {
				fmt.Printf("✅ %s (%s) : OK - %s\n", res.InputTarget.Id, res.InputTarget.Path, res.Err)
			}
		}

		// Exporter les résultats si outputFilePath est spécifié
		if outPath != "" {
			err := reporter.ExportResultsToJsonFile(outPath, finalReport)
			if err != nil {
				fmt.Printf("Erreur lors de l'exportation des résultats: %v\n", err)
			} else {
				fmt.Printf("✅ Résultats exportés vers %s\n", outPath)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "Chemin vers le fichier de configuration JSON (obligatoire)")
	analyzeCmd.Flags().StringVarP(&outPath, "output", "o", "", "Chemin de sortie pour le rapport JSON (optionnel)")
	analyzeCmd.MarkFlagRequired("config")
}