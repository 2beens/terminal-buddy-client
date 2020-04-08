package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var errorUserDataFileCorrupted = errors.New("corrupted user data file")
var errorUserDataEmpty = errors.New("empty user data file")

type ServerSettings struct {
	ServerProtocol string
	ServerAddress  string
	ServerPort     string
}

type Settings struct {
	ServerSettings *ServerSettings

	Username         string
	PasswordHash     string
	SettingsFileName string
}

func NewSettings(serverSettings *ServerSettings, settingsFileName string) (*Settings, error) {
	s := &Settings{
		ServerSettings:   serverSettings,
		SettingsFileName: settingsFileName,
	}

	if err := s.initSettings(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Settings) GetRequestUrl(path string) string {
	if len(s.ServerSettings.ServerPort) > 0 {
		return fmt.Sprintf("%s://%s:%s/%s", s.ServerSettings.ServerProtocol, s.ServerSettings.ServerAddress, s.ServerSettings.ServerPort, path)
	}
	return fmt.Sprintf("%s://%s/%s", s.ServerSettings.ServerProtocol, s.ServerSettings.ServerAddress, path)
}

func (s *Settings) StoreUserData(user *User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	preferencesPath, err := s.preferencesFolder()
	if err != nil {
		return fmt.Errorf("cannot get preferences folder: %w", err)
	}

	settingsData := fmt.Sprintf("%s::%s", user.Username, user.PasswordHash)

	// write the whole body at once
	err = ioutil.WriteFile(preferencesPath+string(os.PathSeparator)+s.SettingsFileName, []byte(settingsData), 0644)
	if err != nil {
		return fmt.Errorf("failed to write settings to file: %w", err)
	}

	return nil
}

func (s *Settings) GetUserData() (*User, error) {
	preferencesPath, err := s.preferencesFolder()
	if err != nil {
		return nil, fmt.Errorf("cannot get preferences folder: %w", err)
	}

	userDataBytes, err := ioutil.ReadFile(preferencesPath + string(os.PathSeparator) + s.SettingsFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings from file: %w", err)
	}

	if len(userDataBytes) == 0 {
		return nil, errorUserDataEmpty
	}

	userData := strings.Split(string(userDataBytes), "::")
	if len(userData) != 2 {
		return nil, errorUserDataFileCorrupted
	}

	return &User{
		Username:     userData[0],
		PasswordHash: userData[1],
	}, nil
}

func (s *Settings) ClearUserData() error {
	preferencesPath, err := s.preferencesFolder()
	if err != nil {
		return fmt.Errorf("cannot get preferences folder: %w", err)
	}

	err = os.Remove(preferencesPath + string(os.PathSeparator) + s.SettingsFileName)
	if err != nil {
		return fmt.Errorf("failed to delte previous settings file: %w", err)
	}

	return s.initSettings()
}

func (s *Settings) preferencesFolder() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	preferencesPath := fmt.Sprintf("%s/Library/Preferences/terminal-buddy", homePath)
	if _, err := os.Stat(preferencesPath); os.IsNotExist(err) {
		//err = os.Mkdir(preferencesPath, os.ModeDir)
		// MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error
		err = os.MkdirAll(preferencesPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return preferencesPath, nil
}

func (s *Settings) initSettings() error {
	preferencesPath, err := s.preferencesFolder()
	if err != nil {
		return fmt.Errorf("cannot get preferences folder: %w", err)
	}

	_, err = os.OpenFile(preferencesPath+string(os.PathSeparator)+s.SettingsFileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("cannot create/open settings file: %w", err)
	}

	return nil
}
