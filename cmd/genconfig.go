package cmd

import (
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/logging"

	"github.com/spf13/cobra"
)

var genconfigCmd = &cobra.Command{
	Use:   "genconfig",
	Short: "Generates a new default config file at the config location.",
	Long:  "Generates a new default config file at the config location. If a config file exists, it will only be overridden when the --reset flag is set to true",
	Run: func(cmd *cobra.Command, args []string) {
		logging.Log.Info.Println("Generating config...")
		reset, _ := cmd.Flags().GetBool("reset")
		if reset {
			err := config.WriteDefaultConfig()
			if err != nil {
				logging.Log.Error.Fatalf("Error writing default config: %s\n", err)
			}
		} else {
			err := config.Read()
			if err != nil {
				logging.Log.Error.Fatalf("Error reading config: %s\n", err)

			}
		}
		logging.Log.Info.Println("Config generated successfully")
	},
}

func init() {
	rootCmd.AddCommand(genconfigCmd)
	genconfigCmd.Flags().Bool("reset", false, "Reset the config file to the default values")

}
