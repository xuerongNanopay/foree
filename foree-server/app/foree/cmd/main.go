package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	foree_boot "xue.io/go-pay/app/foree/cmd/boot"
)

var rootCmd = &cobra.Command{
	Use:   "Foree",
	Short: "Foree Application",
	Long:  `Foree Application`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var appName string
var envFile string
var certPath string

var appCmd = &cobra.Command{
	Use: "app",
	Run: func(cmd *cobra.Command, args []string) {
		envFile, _ := cmd.Flags().GetString("env-file")

		if args[0] == "foree_app" {
			app := &foree_boot.ForeeApp{}
			if err := app.Boot(envFile); err != nil {
				log.Fatal(err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

func main() {
	rootCmd.AddCommand(appCmd)
	appCmd.Flags().String("env-file", "../deploy/.local_env", "Environment file dir.")
	execute()
}
