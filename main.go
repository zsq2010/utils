∂package main

import (
	"log"

	"github.com/zsq2010/utils/notify"
	"gopkg.in/ini.v1"
)

func main() {
	log.Println("--- Starting Notification Test ---")

	// Load configuration from config.ini
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Error: Failed to load config.ini. Make sure the file exists. Details: %v", err)
	}

	// --- Prepare Notifiers ---
	var notifiers []notify.Notifier

	// 1. Configure Bark Notifier
	barkKey := cfg.Section("Bark").Key("Key").String()
	if barkKey != "" {
		barkNotifier := notify.NewBark(notify.BarkConfig{
			Key: barkKey,
		})
		notifiers = append(notifiers, barkNotifier)
		log.Println("Bark notifier configured.")
	} else {
		log.Println("Bark key not found, skipping Bark notifier.")
	}

	// 2. Configure Outlook Notifier
	outlookUser := cfg.Section("Email").Key("outlook_email").String()
	outlookPass := cfg.Section("Email").Key("outlook_password").String()
	if outlookUser != "" && outlookPass != "" {
		outlookNotifier := notify.NewEmail(notify.EmailConfig{
			Provider: notify.Outlook,
			Username: outlookUser,
			Password: outlookPass,
			To:       []string{outlookUser}, // Sending a test email to self
		})
		notifiers = append(notifiers, outlookNotifier)
		log.Println("Outlook notifier configured.")
	} else {
		log.Println("Outlook credentials not found, skipping Outlook notifier.")
	}

	// 3. Configure Gmail Notifier
	gmailUser := cfg.Section("Email").Key("gmail_email").String()
	gmailPass := cfg.Section("Email").Key("gmail_password").String()
	if gmailUser != "" && gmailPass != "" {
		gmailNotifier := notify.NewEmail(notify.EmailConfig{
			Provider: notify.Gmail,
			Username: gmailUser,
			Password: gmailPass,
			To:       []string{gmailUser}, // Sending a test email to self
		})
		notifiers = append(notifiers, gmailNotifier)
		log.Println("Gmail notifier configured.")
	} else {
		log.Println("Gmail credentials not found, skipping Gmail notifier.")
	}

	if len(notifiers) == 0 {
		log.Fatalln("Error: No notifiers were configured. Please check your config.ini.")
	}

	// --- Create a Multi-Channel Notifier ---
	// This will send the message to all configured notifiers in parallel.
	multiNotifier := notify.NewMultiParallel(notifiers...)
	log.Printf("Multi-notifier configured with %d channels. Sending test message...", len(notifiers))

	// --- Send the Message ---
	message := notify.Message{
		Title: "Utils Notification Test",
		Body:  "This is a test message sent from the main.go example.",
	}

	err = multiNotifier.Send(message)
	if err != nil {
		log.Fatalf("--- Test Failed --- \nError sending notifications: %v", err)
	}

	log.Println("--- Test Successful ---")
	log.Println("Notifications sent to all configured channels successfully.")
}
