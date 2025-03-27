package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/inovacc/ksuid"
	"medflow/internal/common"
	"time"

	"github.com/dyammarcano/alfanumeric-cnpj/pkg/cnpj"
	"github.com/nats-io/nats.go"
	_ "modernc.org/sqlite"
)

func ConnectToNATS() (*nats.Conn, nats.JetStreamContext, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, nil, err
	}
	return nc, js, nil
}

func InitSQLite(service string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s.db", service))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return db, nil
}

func SaveEventToSQLite(db *sql.DB, event common.ClinicalEvent) error {
	dataJSON, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	metaJSON, err := json.Marshal(event.Metadata)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO eventos_procesados
		(parent_id, current_id, step, status, timestamp, data, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		event.ParentID, event.CurrentID, event.Step, event.Status, time.Now().Format(time.RFC3339), string(dataJSON), string(metaJSON))
	if err != nil {
		return err
	}
	return nil
}

func SaveEventToPostgres(db *sql.DB, event common.ClinicalEvent) error {
	dataJSON, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	metaJSON, err := json.Marshal(event.Metadata)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO eventos (parent_id, current_id, step, status, timestamp, data, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		event.ParentID, event.CurrentID, event.Step, event.Status, event.Timestamp, dataJSON, metaJSON)
	if err != nil {
		return err
	}
	return nil
}

func GeneratePatient() common.PatientEvent {
	return common.PatientEvent{
		CurrentID: ksuid.NewString(),
		PatientID: uuid.NewString(),
		Step:      "initial_stage",
		Status:    "pending",
		Timestamp: time.Now().Format(time.RFC3339),
		Patient: common.Patient{
			ID:         cnpj.FormatCNPJ(cnpj.GenerateCNPJ()),
			FirstName:  gofakeit.FirstName(),
			SecondName: gofakeit.FirstName(),
			MiddleName: gofakeit.MiddleName(),
			LastName:   gofakeit.LastName(),
			Age:        gofakeit.Number(1, 100),
			Phone:      gofakeit.Phone(),
			Email:      gofakeit.Email(),
			Address:    gofakeit.Address().Address,
		},
	}
}
