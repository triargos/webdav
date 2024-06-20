package cmd

import (
	"github.com/triargos/webdav/pkg/auth"
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/logging"
	"github.com/triargos/webdav/pkg/server"

	"github.com/spf13/cobra"
)

// adduserCmd represents the adduser command
var adduserCmd = &cobra.Command{
	Use:   "adduser",
	Short: "Add a user to the webdav server configuration",
	Run: func(cmd *cobra.Command, args []string) {
		username := cmd.Flag("username").Value.String()
		password := cmd.Flag("password").Value.String()
		admin, _ := cmd.Flags().GetBool("admin")
		dir := cmd.Flag("dir").Value.String()
		jailed, _ := cmd.Flags().GetBool("jailed")
		subdirectories, _ := cmd.Flags().GetStringArray("subdirs")
		config.AddUser(username, config.User{
			Password:       auth.GenHash([]byte(password)),
			Admin:          admin,
			SubDirectories: subdirectories,
			Jail:           jailed,
			Root:           dir,
		})
		err := config.Write()
		if err != nil {
			logging.Log.Error.Fatalf("Error writing config: %s\n", err)
		}
		server.CreateUserDirectories()
		logging.Log.Info.Printf("User %s added successfully. Please restart the service for changes to take effect\n", username)
	},
}

func init() {
	rootCmd.AddCommand(adduserCmd)
	adduserCmd.Flags().StringP("username", "u", "", "Username of the user to add")
	adduserCmd.Flags().StringP("password", "p", "", "Password of the user to add")
	adduserCmd.Flags().BoolP("admin", "a", false, "Is the user an admin")
	adduserCmd.Flags().BoolP("jailed", "j", false, "Is the user jailed")
	adduserCmd.Flags().StringP("dir", "d", "", "Directory of the user to add")
	adduserCmd.Flags().StringArrayP("subdirs", "s", []string{}, "Subdirectories of the user to add")
}
