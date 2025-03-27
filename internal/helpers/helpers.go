package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nats-io/nats.go"
)

func ConnectToNATS() (*nats.Conn, nats.JetStreamContext) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("NATS connection error:", err)
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal("JetStream error:", err)
	}
	return nc, js
}

func InitSQLite(service string) *sql.DB {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s.db", service))
	if err != nil {
		log.Fatal("SQLite open error:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS eventos_procesados (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		parent_id TEXT,
		current_id TEXT,
		step TEXT,
		status TEXT,
		timestamp TEXT,
		data TEXT,
		metadata TEXT
	)`)
	if err != nil {
		log.Fatal("SQLite create table error:", err)
	}

	return db
}

func SaveEventToSQLite(db *sql.DB, event medflow.ClinicalEvent) {
	dataJSON, _ := json.Marshal(event.Data)
	metaJSON, _ := json.Marshal(event.Metadata)

	_, err := db.Exec(`INSERT INTO eventos_procesados
		(parent_id, current_id, step, status, timestamp, data, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		event.ParentID, event.CurrentID, event.Step, event.Status, time.Now().Format(time.RFC3339), string(dataJSON), string(metaJSON))
	if err != nil {
		log.Println("SQLite insert error:", err)
	}
}
