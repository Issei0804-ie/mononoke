package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// https://takaya030.hatenablog.com/entry/2016/09/04/164354
func main() {
	ctx := context.Background()
	_, pwd, _, _ := runtime.Caller(0)

	dir := filepath.Dir(pwd)
	err := godotenv.Load(dir + "/.env")
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	accessToken := os.Getenv("ACCESS_TOKEN")
	refreshToken := os.Getenv("REFRESH_TOKEN")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{"https://mail.google.com/"},
	}

	expiry, _ := time.Parse("2006-01-02", "2017-07-11")
	token := oauth2.Token{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		RefreshToken: refreshToken,
		Expiry:       expiry,
	}

	client := config.Client(ctx, &token)

	service, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Println(err.Error())
		return
	}

	response, err := service.Users.Messages.List("me").Do()
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, message := range response.Messages {
		do, err := service.Users.Messages.Get("me", message.Id).Do()
		if err != nil {
			log.Println(err.Error())
			return
		}

		for _, header := range do.Payload.Headers {
			if header.Name == "Subject" {
				fmt.Println(header.Value)
			}

		}
	}

}
