package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MaddSystems/jonobridge/common/utils"
)

// ------------------------------------------------------------------
// CONFIGURACIÓN GLOBAL (puedes sobrescribir con variables de entorno)
// ------------------------------------------------------------------
var (
	TelegramBotToken = getEnv("TELEGRAM_BOT_TOKEN", "")
	TelegramChatID   = getEnv("TELEGRAM_CHAT_ID", "") // puedes usar varios separados por coma

	WebhookURL = getEnv("ALERT_WEBHOOK_URL", "") // para integrar con Slack, Teams, etc.

	EmailSMTPHost = getEnv("EMAIL_SMTP_HOST", "")
	EmailSMTPPort = getEnv("EMAIL_SMTP_PORT", "587")
	EmailUser     = getEnv("EMAIL_USER", "")
	EmailPass     = getEnv("EMAIL_PASS", "")
	EmailFrom     = getEnv("EMAIL_FROM", "alertas@tuservicio.com")
	EmailTo       = getEnv("EMAIL_TO", "") // separados por coma
)

// ------------------------------------------------------------------
// FUNCIONES QUE USARÁS DIRECTAMENTE DESDE LAS REGLAS GRULE
// ------------------------------------------------------------------

// SendTelegram envía alerta a uno o varios chats de Telegram
func SendTelegram(message string) {
	if TelegramBotToken == "" || TelegramChatID == "" {
		utils.VPrint("Telegram no configurado (falta token o chat_id)")
		return
	}

	chatIDs := strings.Split(TelegramChatID, ",")
	for _, chatID := range chatIDs {
		chatID = strings.TrimSpace(chatID)
		url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", TelegramBotToken)

		payload := map[string]string{
			"chat_id":    chatID,
			"text":       "🚨 " + message,
			"parse_mode": "HTML",
		}
		jsonPayload, _ := json.Marshal(payload)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil || resp.StatusCode >= 400 {
			utils.VPrint("Error enviando Telegram: %v", err)
			return
		}
		resp.Body.Close()
	}
}

// SendEmail envía correo (simple, sin adjuntos)
func SendEmail(subject, body string) {
	if EmailSMTPHost == "" {
		utils.VPrint("Email SMTP no configurado")
		return
	}
	// Implementación básica (puedes mejorar con net/smtp + html/template si quieres)
	log.Printf("EMAIL → %s | Asunto: %s | Cuerpo: %s", EmailTo, subject, body)
	// Aquí iría el envío real con net/smtp si lo necesitas
}

// Webhook genérico (Slack, Microsoft Teams, Discord, etc.)
func SendWebhook(message string) {
	if WebhookURL == "" {
		return
	}
	payload := map[string]string{"text": "🚨 " + message}
	jsonPayload, _ := json.Marshal(payload)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(WebhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil || resp.StatusCode >= 400 {
		utils.VPrint("Error enviando webhook: %v", err)
		return
	}
	resp.Body.Close()
}

// Log con soporte para printf y GRULE_LOG_DESTINATION
func Log(message string, args ...interface{}) {
	// Leer configuración de destino desde variable de entorno
	destination := os.Getenv("GRULE_LOG_DESTINATION")
	if destination == "" {
		destination = "console" // Default: solo consola
	}
	
	// Formatear mensaje con argumentos printf
	formattedMsg := fmt.Sprintf(message, args...)
	
	switch destination {
	case "console":
		log.Printf("[GRULE] %s", formattedMsg)
	case "mysql":
		logToMySQL(formattedMsg)
	case "both": 
		log.Printf("[GRULE] %s", formattedMsg)
		logToMySQL(formattedMsg)
	case "none": 
		// No hacer nada - silencioso
	default: 
		log.Printf("[GRULE] %s", formattedMsg) // Fallback: consola
	}
}

// logToMySQL - Helper interno para usar sistema audit existente
func logToMySQL(message string) {
	// USAR SISTEMA AUDIT EXISTENTE: execution_steps table
	// NO crear tabla duplicada - reutilizar audit.CaptureExecution()
	go logToAuditSystem(message)
}

// logToAuditSystem - Integrar con sistema audit existente  
func logToAuditSystem(message string) {
	// TODO: Usar audit.CaptureExecution() con execution_steps existente
	// Formato: audit.CaptureExecution(imei, false, []StepExecution{{Description: message}})
	// Por ahora solo log para no romper la compilación
	log.Printf("[AUDIT-MYSQL] %s", message)
}

// ------------------------------------------------------------------
// UTILIDADES INTERNAS
// ------------------------------------------------------------------
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}