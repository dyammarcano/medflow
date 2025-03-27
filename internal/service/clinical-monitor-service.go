package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"

	"medflow/internal/common"
	"medflow/internal/helpers"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan common.ClinicalEvent)
var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ClinicalMonitorService(cmd *cobra.Command, _ []string) error {
	database, err := sql.Open("postgres", "postgres://user:password@localhost:5432/medflow?sslmode=disable")
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}
	defer func(database *sql.DB) {
		_ = database.Close()
	}(database)

	if err := database.PingContext(cmd.Context()); err != nil {
		return fmt.Errorf("database ping error: %w", err)
	}

	if err := createTable(database); err != nil {
		return fmt.Errorf("create table error: %w", err)
	}

	natsConn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("nats connection error: %w", err)
	}
	defer natsConn.Close()

	jetStream, err := natsConn.JetStream()
	if err != nil {
		return fmt.Errorf("jetstream error: %w", err)
	}

	// Ensure stream exists
	_, err = jetStream.AddStream(&nats.StreamConfig{
		Name:     "OPERATION_STREAM",
		Subjects: []string{"operation.*.data"},
		Storage:  nats.FileStorage,
	})
	if err != nil && !strings.Contains(err.Error(), "stream name already in use") {
		return fmt.Errorf("error creating stream: %w", err)
	}

	_, err = jetStream.Subscribe("operation.*.data", func(msg *nats.Msg) {
		var event common.ClinicalEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Println("Invalid JSON:", err)
			return
		}

		event.Timestamp = time.Now().Format(time.RFC3339)
		if err := helpers.SaveEventToPostgres(database, event); err != nil {
			log.Println("Save event error:", err)
			return
		}

		broadcast <- event
	}, nats.Durable("monitor-durable"), nats.ManualAck())
	if err != nil {
		return fmt.Errorf("subscription error: %w", err)
	}

	http.HandleFunc("/ws", handleWebSocket)
	go handleMessages()

	log.Println("Monitor listening on ws://localhost:8080/ws")
	return http.ListenAndServe(":8080", nil)
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS eventos (
		id SERIAL PRIMARY KEY,
		parent_id TEXT,
		current_id TEXT,
		step TEXT,
		status TEXT,
		timestamp TEXT,
		data JSONB,
		metadata JSONB
	)`)
	return err
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	clients[conn] = true
	log.Println("Client connected")
}

func handleMessages() {
	for {
		event := <-broadcast
		message, _ := json.Marshal(event)
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("WebSocket write error:", err)
				_ = client.Close()
				delete(clients, client)
			}
		}
	}
}
