package cmd

import (
	"fmt"
	"os"

	"github.com/2beens/term-buddy-commander/internal"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var settings *internal.Settings
var loggedUser *internal.User

var useLocalServer = false

// TODO: make configurable, via params, or some config file, or something else
const (
	serverLocalProtocol = "http"
	serverLocalAddress  = "localhost"
	serverLocalPort     = "8080"
	serverProtocol      = "https"
	serverAddress       = "www.serjspends.de/tb"
	serverPort          = ""

	settingsFilename       = "term-buddy-settings"
	passwordHashHeaderName = "Term-Buddy-Pass-Hash"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "term-buddy-commander",
	Short: "Your terminal buddy, helping you with reminders and notes",
	Long: `Your terminal buddy, helping you with reminders and notes
		- write reminders with: remind
		- write notes with: note
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("in Run() of root cmd")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	fmt.Println("in Execute() of root cmd")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func SetLoggedUser(username, passwordHash string) {
	loggedUser = &internal.User{
		Username:     username,
		PasswordHash: passwordHash,
	}
	if err := settings.StoreUserData(loggedUser); err != nil {
		log.Errorf("failed to store user data: %s", err.Error())
	}
}

func UserLogged() bool {
	return loggedUser != nil
}

func GetServerSettings() *internal.ServerSettings {
	if useLocalServer {
		return &internal.ServerSettings{
			ServerProtocol: serverLocalProtocol,
			ServerAddress:  serverLocalAddress,
			ServerPort:     serverLocalPort,
		}
	}
	return &internal.ServerSettings{
		ServerProtocol: serverProtocol,
		ServerAddress:  serverAddress,
		ServerPort:     serverPort,
	}
}

func init() {
	fmt.Println("in init() of root cmd")

	// TODO: log to terminal for now
	log.SetOutput(os.Stdout)

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.term-buddy-commander.yaml)")
	rootCmd.PersistentFlags().BoolVar(&useLocalServer, "ls", false, "Set to 'true' to connect to localhost server")

	// Cobra also supports local flags, which will only run when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	fmt.Println("in initConfig() of root cmd")

	if useLocalServer {
		log.Warn("using local server")
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".term-buddy-commander" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".term-buddy-commander")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	var err error
	settings, err = internal.NewSettings(GetServerSettings(), settingsFilename)
	if err != nil {
		panic(err)
	}

	loggedUser, err = settings.GetUserData()
	if err != nil {
		log.Warn("user data empty/corrupted, please login/register")
	}
}
