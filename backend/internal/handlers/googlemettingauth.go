package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/meet/v2"
	"google.golang.org/api/option"
)

// getClient uses a Context and Config to retrieve a Token then generate a Client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "meet_oauth.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// CreateSpace creates a new Google Meet space and returns the URI.
func CreateSpace() (string, error) {
	ctx := context.Background()

	// 1. Read the credentials file
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return "", fmt.Errorf("unable to read client secret file: %v", err)
	}

	// 2. Configure the scope
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/meetings.space.created")
	if err != nil {
		return "", fmt.Errorf("unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	// 3. Create the Meet Service
	srv, err := meet.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("unable to retrieve Meet client: %v", err)
	}

	// 4. Create the Space request
	space := &meet.Space{}
	res, err := srv.Spaces.Create(space).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create space: %v", err)
	}

	fmt.Printf("Meet URL: %s\n", res.MeetingUri)
	return res.MeetingUri, nil
}

// --- Helper Functions for OAuth2 Token Management ---
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Force the refresh token and ensure the redirect matches your Google Console
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	fmt.Printf("1. Open this link in your browser: \n%v\n\n", authURL)
	fmt.Printf("2. After login, you will be redirected to localhost:8080/dashboard.\n")
	fmt.Printf("3. Look at the URL in your browser address bar. Copy the 'code=' value.\n")
	fmt.Printf("4. Paste the code here: ")

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
