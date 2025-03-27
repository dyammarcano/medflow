package service

import (
	"database/sql"
	"encoding/json"
	"github.com/spf13/cobra"
	"log"
	"medflow/internal/common"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan common.ClinicalEvent)
var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ClinicalMonitorService(cmd *cobra.Command, args []string) error {
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/medflow?sslmode=disable")
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Fatal("Close error:", err)
		}
	}(db)

	if err := createTable(db); err != nil {
		return err
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return err
	}

	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	_, err = js.Subscribe("operation.*.data", func(msg *nats.Msg) {
		var event common.ClinicalEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Println("Invalid JSON:", err)
			return
		}

		event.Timestamp = time.Now().Format(time.RFC3339)
		saveEventToPostgres(db, event)
		broadcast <- event
	}, nats.Durable("monitor-durable"), nats.ManualAck())
	if err != nil {
		return err
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
	if err != nil {
		return err
	}
	return nil
}

func saveEventToPostgres(db *sql.DB, event common.ClinicalEvent) {
	dataJSON, err := json.Marshal(event.Data)
	if err != nil {
		log.Println("JSON marshal error:", err)
	}

	metaJSON, err := json.Marshal(event.Metadata)
	if err != nil {
		log.Println("JSON marshal error:", err)
	}

	_, err = db.Exec(`INSERT INTO eventos (parent_id, current_id, step, status, timestamp, data, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		event.ParentID, event.CurrentID, event.Step, event.Status, event.Timestamp, dataJSON, metaJSON)
	if err != nil {
		log.Println("Insert error:", err)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	clients[conn] = true
}

func handleMessages() {
	for {
		event := <-broadcast
		message, _ := json.Marshal(event)
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("WebSocket write error:", err)
				if err := client.Close(); err != nil {
					log.Println("WebSocket close error:", err)
					return
				}
				delete(clients, client)
			}
		}
	}
}
