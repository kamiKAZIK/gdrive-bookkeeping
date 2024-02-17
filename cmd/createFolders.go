package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createFoldersCmd represents the createFolders command
var createFoldersCmd = &cobra.Command{
	Use:   "createFolders",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("createFolders called")
	},
}

func init() {
	rootCmd.AddCommand(createFoldersCmd)

	createFoldersCmd.Flags().IntP("year", "y", time.Now().Year(), "Year for which the folders need to be created")
	createFoldersCmd.MarkFlagRequired("year")

	viper.BindPFlag("year", createFoldersCmd.Flags().Lookup("year"))

// 	t := time.Now()
// 	viper.SetDefault("year", t.Year())
}
