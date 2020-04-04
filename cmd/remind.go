package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/2beens/term-buddy-commander/internal"
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

		if len(args) > 0 && args[0] == "all" {
			handleAll()
			return
		}

		fmt.Printf("args (%d): %v", len(args), args)
	},
}

func handleAll() {
	client := http.Client{}
	url := fmt.Sprintf("%s://%s:%s/remind/%s/all", serverProtocol, serverAddress, serverPort, loggedUser.Username)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set(passwordHashHeaderName, loggedUser.PasswordHash)

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("FIX: server resp status: %d", resp.StatusCode)

	var serverResp internal.ServerResponse
	err = json.NewDecoder(resp.Body).Decode(&serverResp)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(serverResp.Message)
	if !serverResp.Ok {
		log.Println("not OK :(")
		return
	}

	log.Println(serverResp.Data)
	if reminders, ok := serverResp.Data.([]internal.Reminder); ok {
		log.Printf("%v", reminders)
	} else {
		log.Error("cannot cast server response data to []Reminder")
	}
}

func init() {
	fmt.Println("in init() of remind cmd")

	rootCmd.AddCommand(remindCmd)

}
