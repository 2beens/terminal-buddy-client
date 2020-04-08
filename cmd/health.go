package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/2beens/term-buddy-commander/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("calling health...")
		log.Printf("logged in: %t", UserLogged())
		reqUrl := settings.GetRequestUrl("health")
		log.Warnf("req url: %s", reqUrl)

		client := http.Client{}
		request, err := http.NewRequest("GET", reqUrl, nil)
		if err != nil {
			log.Fatalln(err)
		}

		resp, err := client.Do(request)
		if err != nil {
			log.Fatalln(err)
		}

		var serverResp internal.ServerResponse
		err = json.NewDecoder(resp.Body).Decode(&serverResp)
		if err != nil {
			log.Fatalln(err)
		}

		// TODO: response from server is not logged, but displayed to user
		if serverResp.Ok {
			log.Println(serverResp.Message)
		} else {
			log.Println("error: " + serverResp.Message)
		}
	},
}

func init() {
	fmt.Println("in init() of health cmd")
	rootCmd.AddCommand(healthCmd)
}
