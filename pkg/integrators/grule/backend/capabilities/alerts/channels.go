package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func (c *AlertsCapability) Log(message string) {
	log.Printf("[GRULE] %s", message)
}

func (c *AlertsCapability) SendTelegram(message string) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	chatIDs := os.Getenv("TELEGRAM_CHAT_ID")

	if token == "" || chatIDs == "" {
		// Only log if verbose or just once, but standard log is fine
		// log.Println("Telegram not configured") 
		return
	}

	for _, chatID := range strings.Split(chatIDs, ",") {
		chatID = strings.TrimSpace(chatID)
		url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

		payload := map[string]string{
			"chat_id":    chatID,
			"text":       "🚨 " + message,
			"parse_mode": "HTML",
		}
		jsonPayload, _ := json.Marshal(payload)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil || resp.StatusCode >= 400 {
			log.Printf("Error sending Telegram to %s: %v", chatID, err)
			continue
		}
		resp.Body.Close()
	}
}
