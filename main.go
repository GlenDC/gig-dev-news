package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	// fetch token from a token file or directly from the command line
	token := fetchToken()

	// create the Telegram bot, using the given token to authenticate
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(fmt.Errorf(
			"couldn't authenticate/create dev news Telegram bot: %v", err))
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func fetchToken() string {
	if len(os.Args) > 2 {
		log.Fatal("too many arguments")
	}

	// if no arguments are given, try to read the default path
	if len(os.Args) == 1 {
		log.Println("no positional argument given")
		token, err := readTokenFile(".token")
		if err != nil {
			log.Fatalf("couldn't read default ./.token file: %v", err)
		}
		return token
	}
	arg := os.Args[1]

	// first try to use the argument as a file name
	token, err := readTokenFile(arg)
	if err != nil {
		if err == os.ErrNotExist {
			log.Println("token given as positional argument")
			return token
		}
	}

	log.Println("token read from token file " + arg)
	return token
}

func readTokenFile(path string) (string, error) {
	tokenFile, err := os.Open(path)
	if err != nil {
		return "", err
	}

	// we could read the token file, let's assume the first line is the token
	var token string
	_, err = fmt.Fscanln(tokenFile, &token)
	if err != nil {
		return "", fmt.Errorf("couldn't read token (first line) in token file %s: %v", path, err)
	}
	return token, nil
}
