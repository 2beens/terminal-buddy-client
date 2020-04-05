package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/2beens/term-buddy-commander/internal"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")

		if len(args) == 1 && args[0] == "clear" {
			handleClearUserData()
			return
		}

		handleLogin(args)
	},
}

func handleLogin(args []string) {
	if len(args) != 2 {
		log.Error("no, no. wrong args. 1st should be username, 2nd password please")
		return
	}

	username := args[0]
	password := args[1]
	passwordHash := HashPassword(password)

	data := url.Values{}
	data.Add("username", username)
	data.Add("password_hash", passwordHash)

	client := http.Client{}
	request, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s://%s:%s/user/login", serverProtocol, serverAddress, serverPort),
		bytes.NewBufferString(data.Encode()),
	)
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

	if serverResp.Ok {
		log.Println(serverResp.Message)
		SetLoggedUser(username, passwordHash)
	} else {
		log.Println("error: " + serverResp.Message)
	}

	// TODO: user is returned in serverResp.Data
}

func handleClearUserData() {
	log.Warn("clearing log data ...")
	if err := settings.ClearUserData(); err != nil {
		log.Errorf("user data not cleared: %s", err.Error())
	} else {
		log.Debug("user data cleared")
	}
}

func init() {
	fmt.Println("in init() of login cmd")
	rootCmd.AddCommand(loginCmd)
}
