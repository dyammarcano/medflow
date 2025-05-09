package helpers

import (
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
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
		data JSONB
	)`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func SaveEventToSQLite(db *sql.DB, event common.PatientEvent) error {
	_, err := db.Exec(`INSERT INTO eventos_procesados
		(parent_id, current_id, step, status, timestamp, data)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		event.PatientID, event.CurrentID, event.Step, event.Status, time.Now().Format(time.RFC3339), event.String())
	if err != nil {
		return err
	}
	return nil
}

func SaveEventToPostgres(db *sql.DB, event common.PatientEvent) error {
	_, err := db.Exec(`INSERT INTO eventos (parent_id, current_id, step, status, timestamp, data)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		event.PatientID, event.CurrentID, event.Step, event.Status, event.Timestamp, event.String())
	if err != nil {
		return err
	}
	return nil
}

func GeneratePatient() common.PatientEvent {
	return common.PatientEvent{
		CurrentID: ksuid.NewString(),
		PatientID: ksuid.NewString(),
		Step:      "patient-intake",
		Status:    "initiated",
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
