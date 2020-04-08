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

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("register called with args: %v\n", args)
		if len(args) != 2 {
			fmt.Println("no, no. wrong args. 1st should be username, 2nd password please")
			return
		}

		username := args[0]
		password := args[1]
		passwordHash := HashPassword(password)

		reqUrl := settings.GetRequestUrl("user/register")
		log.Warnf("req url: %s", reqUrl)

		data := url.Values{}
		data.Set("username", username)
		data.Add("password_hash", passwordHash)

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

		if serverResp.Ok {
			log.Println(serverResp.Message)
			SetLoggedUser(username, passwordHash)
		} else {
			log.Println("error: " + serverResp.Message)
		}
	},
}

func init() {
	fmt.Println("in init() of register cmd")

	rootCmd.AddCommand(registerCmd)

}
