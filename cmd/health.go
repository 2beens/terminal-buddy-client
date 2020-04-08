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
		reqUrl := getRequestUrl("health")
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

		//respBytes, _ := ioutil.ReadAll(resp.Body)
		//log.Println(string(respBytes))

		//var result map[string]interface{}
		var serverResp internal.ServerResponse
		err = json.NewDecoder(resp.Body).Decode(&serverResp)
		if err != nil {
			log.Fatalln(err)
		}

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// healthCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// healthCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
