package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/2beens/term-buddy-commander/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	regexpRemindInCommand = regexp.MustCompile(`^(.+)\sin\s(\w+)$`)
	regexpRemindAtCommand = regexp.MustCompile(`^(.+)\sat\s(\w+)$`)
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

		if len(args) == 0 {
			log.Error("no arguments specified")
			return
		}

		// TODO: use subcommands instead
		if len(args) == 1 && args[0] == "all" {
			handleAll()
			return
		}

		// TODO: use subcommands instead
		if len(args) == 2 && args[0] == "all" && args[1] == "clear" {
			handleClearAll()
			return
		}

		command := strings.Join(args, " ")
		log.Printf("args (%d): %v\n", len(args), args)
		log.Printf("command: %s", command)

		remindInMatched := regexpRemindInCommand.MatchString(command)
		remindAtMatched := regexpRemindAtCommand.MatchString(command)
		if !remindInMatched && !remindAtMatched {
			log.Error("wrong remind message format: " + command)
			return
		}

		if remindInMatched {
			groups := regexpRemindInCommand.FindStringSubmatch(command)
			handleNewRemind(groups[1], groups[2])
		} else {
			// TODO: support for specific time
			// 		e.g. remind "cheeki breeki" at 5pm
			log.Warn("'at' not supported yet, sorry")
		}
	},
}

// remind "kick the dog" in "3 hours"
// duration instructions, such as "300ms", "-1.5h" or "2h45m".
func handleNewRemind(remindMessage, durationInstructions string) {
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

	// TODO: response from server is not logged, but displayed to user
	log.Println("server:")
	if !serverResp.Ok {
		log.Println("not OK :(")
	}
	log.Println("\t- " + serverResp.Message)
	if len(serverResp.DataJsonBytes) > 0 {
		log.Printf("\t- %v", string(serverResp.DataJsonBytes))
	}
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

func handleClearAll() {
	log.Warn("not supported yet")
	return
}

func init() {
	fmt.Println("in init() of remind cmd")
	rootCmd.AddCommand(remindCmd)
}
