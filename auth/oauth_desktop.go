package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

func NewDeskTopAccountClient(configDir, desktopAccountFile string) (*http.Client, error) {
	desktopAccountFilePath := filepath.Join(configDir, desktopAccountFile)
	b, err := os.ReadFile(desktopAccountFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	return getClient(configDir, config), nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(configDir string, config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tokFilePath := filepath.Join(configDir, tokFile)
	tok, err := tokenFromFile(tokFilePath)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFilePath, tok)
	}
	if !tok.Valid() {
		tok = refreshToken(config, tok)
		saveToken(tokFilePath, tok)
	}

	return config.Client(context.Background(), tok)
}

func refreshToken(config *oauth2.Config, token *oauth2.Token) *oauth2.Token {
	tok, err := config.TokenSource(context.Background(), token).Token()
	if err != nil {
		log.Fatalf("Unable to refresh token %v", err)
	}
	return tok
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer silentClose(f)
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer silentClose(f)
	silentError(json.NewEncoder(f).Encode(token))
}
