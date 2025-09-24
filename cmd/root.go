package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "loganizer",
	Short: "Outil d'analyse de logs avec des règles personnalisées",
	Long: `loganizer est un outil en ligne de commande pour analyser des fichiers de logs`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
		os.Exit(1)
	}
}

func init() {}