package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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

		fmt.Printf("args (%d): %v\n", len(args), args)

		if len(args) == 1 && args[0] == "all" {
			handleAll()
			return
		}

		// TODO: do this in a smarter way :)
		if len(args) == 3 && strings.ToLower(args[1]) == "in" {
			handleNewRemind(args)
		}
	},
}

// remind "kick the dog" in "3 hours"
func handleNewRemind(args []string) {
	remindMessage := args[0]
	// such as "300ms", "-1.5h" or "2h45m".
	durationInstructions := args[2]
	duration, err := time.ParseDuration(durationInstructions)
	if err != nil {
		log.Errorf("provided duration is invalid: %s", durationInstructions)
		return
	}

	dueDate := time.Now().Add(duration)
	dueDateUnixTimestamp := dueDate.Unix()

	log.Printf("due date param: %+v", dueDate)

	reqUrl := settings.GetRequestUrl(fmt.Sprintf("remind/%s", loggedUser.Username))
	log.Warnf("req url: %s", reqUrl)

	data := url.Values{}
	data.Add("message", remindMessage)
	data.Add("due_date", strconv.FormatInt(dueDateUnixTimestamp, 10))
	data.Add("password_hash", loggedUser.PasswordHash)

	client := http.Client{}
	request, err := http.NewRequest("POST", reqUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("server resp status: %d", resp.StatusCode)

	var serverResp internal.ServerResponse
	err = json.NewDecoder(resp.Body).Decode(&serverResp)
	if err != nil {
		log.Fatalln(err)
	}

	if !serverResp.Ok {
		log.Error("error: " + serverResp.Message)
		return
	}

	log.Println("server:")
	if !serverResp.Ok {
		log.Println("not OK :(")
	}
	log.Println("\t- " + serverResp.Message)
	log.Printf("\t- %v", string(serverResp.DataJsonBytes))
}

func handleAll() {
	reqUrl := settings.GetRequestUrl(fmt.Sprintf("remind/%s/all", loggedUser.Username))
	log.Warnf("req url: %s", reqUrl)

	client := http.Client{}
	request, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}
	request.Header.Set(passwordHashHeaderName, loggedUser.PasswordHash)

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("resp status: %d", resp.StatusCode)

	var serverResp internal.ServerResponse
	err = json.NewDecoder(resp.Body).Decode(&serverResp)
	if err != nil {
		log.Fatalln(err)
	}

	if !serverResp.Ok {
		log.Println("\t- " + serverResp.Message)
		return
	}

	reminders := &[]internal.Reminder{}
	err = json.Unmarshal(serverResp.DataJsonBytes, reminders)
	if err != nil {
		log.Errorf("reminders unmarshal error: %s", err.Error())
		return
	}

	if len(*reminders) == 0 {
		log.Printf("no reminders yet. go on, make one")
		return
	}

	for _, r := range *reminders {
		log.Printf(" - %d: %s", r.Id, r.Message)
		log.Printf("\t- %v", time.Unix(r.DueDate, 0))
	}
}

func init() {
	fmt.Println("in init() of remind cmd")

	rootCmd.AddCommand(remindCmd)

}
