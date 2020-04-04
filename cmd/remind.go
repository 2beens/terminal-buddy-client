package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var remindCmd = &cobra.Command{
	Use:   "remind",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("remind called")

		if !UserLogged() {
			log.Warn("not logged in")
			return
		}

		fmt.Printf("args (%d): %v", len(args), args)
	},
}

func init() {
	fmt.Println("in init() of remind cmd")

	rootCmd.AddCommand(remindCmd)

}
