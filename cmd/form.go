package cmd

import (
	"github.com/spf13/cobra"
)

// formCmd represents the form command
var formCmd = &cobra.Command{
	Use:   "form",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(formCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// formCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// formCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Field to descript field of a Form
type Field struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
	Array    bool   `json:"array"`
}
