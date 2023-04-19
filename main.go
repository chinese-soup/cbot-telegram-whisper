package main

import (
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"os"
)

func main() {
	apiToken, ok := os.LookupEnv("TELEGRAM_APITOKEN")
	if !ok {
		panic("No TELEGRAM_APITOKEN environment variable found.")
	}

	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		panic(err)
	}

	whisperModel := GetModel()

	bot.Debug = true
	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}

		if update.Message.Voice == nil {
			continue
		}

		file, err := bot.GetFileDirectURL(update.Message.Voice.FileID)
		if err != nil {
			log.Fatal(err)
		}

		//TODO: Don't save to a temporary file, replace with a pipe to ffmpeg
		resp, err := grab.Get("tmp/", file)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Download saved to", resp.Filename)

		voiceMsgFilename := resp.Filename
		err = ffmpeg.Input(voiceMsgFilename).
			Output("tmp/tmp.wav",
				ffmpeg.KwArgs{"acodec": "pcm_s16le", "ac": "1", "ar": "16000"}).
			OverWriteOutput().ErrorToStdOut().Run()

		if err != nil {
			log.Fatal(err)
		}

		samples, err := GetSamplesFromFilePath("tmp/tmp.wav")
		if err != nil {
			log.Printf("Unfortunately this happened: %v", err)
			continue
		}
		recognizedText := ProcessSamples(whisperModel, samples)

		// Now that we know we've gotten a new message, we can construct a
		// reply! We'll take the Chat ID and Text from the incoming message
		// and use it to create a new message.
		if recognizedText != "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, recognizedText)
			// We'll also say that this message is a reply to the previous message.
			// For any other specifications than Chat ID or Text, you'll need to
			// set fields on the `MessageConfig`.
			msg.ReplyToMessageID = update.Message.MessageID

			// Okay, we're sending our message off! We don't care about the message
			// we just sent, so we'll discard it.
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending message: %v", err) //TODO: Handle and retry?
			}
		} else {
			log.Println("Nothing recognized, nothing to send.")
		}
	}
}
