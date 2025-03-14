package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/kenliu/peer-acks-v2/app/slack"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	functions.HTTP("SlackEvents", HandleSlackEvents)
	functions.HTTP("SlackSlashCommand", HandleSlackSlashCommand)

	// Initialize DB connection
	datasource := os.Getenv("DATASOURCE")
	log.Printf("Connecting to database with datasource: %s", datasource)
	var err error
	db, err = sql.Open("postgres", datasource)
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	log.Println("Successfully connected to database")
}

func main() {
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	// By default, listen on all interfaces. If testing locally, run with
	// LOCAL_ONLY=true to avoid triggering firewall warnings
	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	}

	if err := funcframework.StartHostPort(hostname, port); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}

// HandleSlackEvents handles POST /slack/events
func HandleSlackEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	challenge, err := slack.HandleChallengeEvent(body)
	if err != nil {
		http.Error(w, "Failed to handle Slack event", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(challenge))
}

// HandleSlackSlashCommand handles POST /slack/slashcommand
func HandleSlackSlashCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	verr := slack.ValidateRequestSignature(r.Header, body, os.Getenv("SLACK_SIGNING_SECRET"))
	if verr != nil {
		log.Println(verr)
		log.Println("request signature verification failed")
		http.Error(w, "Invalid request signature", http.StatusForbidden)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	userName := r.FormValue("user_name")
	message := r.FormValue("text")
	userId := r.FormValue("user_id")
	log.Println("received slack slash command from user:" + userName + " with message text: \"" + message + "\"")

	responseMessage, err := slack.HandleSlashCommand(message, userId, db)
	if err != nil {
		http.Error(w, "Failed to handle slash command", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(responseMessage))
}
