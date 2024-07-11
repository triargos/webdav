package cmd

import (
	"github.com/triargos/webdav/pkg/config"
	"github.com/triargos/webdav/pkg/environment"
	"github.com/triargos/webdav/pkg/fs"
	"github.com/triargos/webdav/pkg/user"
	"log/slog"
	"os"

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
		configService := config.NewViperConfigService(environment.NewOsEnvironmentService())
		fsService := fs.NewOsFileSystemService()
		userService := user.NewOsUserService(configService, fsService)
		addUserErr := userService.AddUser(username, config.User{
			Password:       password,
			Admin:          admin,
			SubDirectories: subdirectories,
			Jail:           jailed,
			Root:           dir,
		})
		if addUserErr != nil {
			slog.Error("failed to add user", "error", addUserErr.Error())
			os.Exit(1)
		}
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
