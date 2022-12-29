package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

var (
	userID       = os.Getenv("USERID")
	channelID    = os.Getenv("CHANNELID")
	webhookName  = "webhook"
	botToken     = os.Getenv("BOTTOKEN")     // developer portal -> app -> bot -> token
	webhookToken = os.Getenv("WEBHOOKTOKEN") // app -> edit channel -> integrations -> webhooks -> copywebhook URL
)

func main() {
	logger := zerolog.New(os.Stdout)

	err := ReadConfig()
	if err != nil {
		logger.Error().Err(err).Msg("error occurred during read config")
		return
	}

	err = BasicBotFlow()
	if err != nil {
		logger.Error().Err(err).Msg("error occurred during basic bot flow")
		return
	}
	logger.Info().Msg("Bot is Running! ")
	// use this to wait message from user
	<-make(chan struct{})

	//  Developer portal -> Bot -> Token, this is how to get the token
	session := createSessionWithToken(logger, botToken)

	// create new webhook and send with that a message
	// there is no need for this, in this way it is necessary to do OAuth2 and give permission to the bot
	hook, err := createWebhook(session, userID, webhookName)
	if err != nil {
		logger.Error().Err(err).Msg("error occurred during create webhook")
		return
	}

	logger.Info().Msg(hook.ChannelID + " token=> " + hook.Token)

	// It is easier to create a webhook from the panel on the discord application and get the token and id.
	sendMessageWithWebhook(session, logger, channelID, "hello from golang", webhookToken)
}

func sendMessageWithWebhook(s *discordgo.Session, logger zerolog.Logger, webhookID, content, token string) {
	msg, err := s.WebhookExecute(
		webhookID,
		token,
		true,
		&discordgo.WebhookParams{Content: content},
	)
	if err != nil {
		logger.Error().Err(err).Msg("error occurred webhook execute")
	}
	logger.Info().Msg(msg.ChannelID + " content=> " + msg.Content)
}

func createSessionWithToken(logger zerolog.Logger, token string) *discordgo.Session {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error().Err(err).Msg("error occurred during create new discord bot")
	}
	return s
}

func createWebhook(s *discordgo.Session, channelID string, name string) (*discordgo.Webhook, error) {
	st, err := s.WebhookCreate(channelID, name, "")
	if err != nil {
		return nil, fmt.Errorf("webhook create error: %v", err)
	}
	return st, nil
}

var BotID string

func BasicBotFlow() error {
	goBot, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		return fmt.Errorf("error occurred during create new discord bot: %v ", err)
	}
	// Making our bot a user using User function .
	u, err := goBot.User("@me")
	if err != nil {
		return fmt.Errorf("error occurred go bot user: %v", err)
	}
	// Storing our id from u to BotId .
	BotID = u.ID

	goBot.AddHandler(sendUserMessage)

	err = goBot.Open()
	if err != nil {
		return fmt.Errorf("error occurred go bot open: %v", err)
	}
	return nil
}

func sendUserMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if strings.Contains(m.Content, "okan") {
		c, err := s.UserChannelCreate("927509052014092318")
		if err != nil {
			fmt.Println("create errror", err)
		}
		content := fmt.Sprintf("this task for you =>%v", m.Content)
		_, err = s.ChannelMessageSend(c.ID, content)
		if err != nil {
			fmt.Println("message send error", err)
		}
	}
}

var (
	Token     string
	BotPrefix string
	conf      *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
	file, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(file))

	err = json.Unmarshal(file, &conf)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// After storing value in config variable we will access it and storing it in our declared variables .
	Token = conf.Token
	BotPrefix = conf.BotPrefix
	return nil
}
