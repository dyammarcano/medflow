package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nats-io/nats.go"

	"medflow/internal/common"
	"medflow/internal/helpers"
)

func StartPatientIntakeService(cmd *cobra.Command, _ []string) error {
	database, err := helpers.InitSQLite("patient_intake")
	if err != nil {
		return fmt.Errorf("database initialization error: %w", err)
	}
	defer func(database *sql.DB) {
		_ = database.Close()
	}(database)

	natsConn, js, err := helpers.ConnectToNATS()
	if err != nil {
		return fmt.Errorf("nats connection error: %w", err)
	}
	defer natsConn.Close()

	// Ensure stream exists
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "OPERATION_STREAM",
		Subjects: []string{common.SubjectOperationIncomingData},
		Storage:  nats.FileStorage,
	})
	if err != nil && !strings.Contains(err.Error(), "stream name already in use") {
		return fmt.Errorf("stream creation error: %w", err)
	}

	sub, err := js.Subscribe(common.SubjectOperationIncomingData, func(msg *nats.Msg) {
		var event = helpers.GeneratePatient()
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Println("Failed to parse message:", err)
			return
		}

		log.Printf("[patient-intake] Processing patient %s", event.PatientID)
		event.Step = "patient-intake"
		event.Status = "completed"
		event.Timestamp = time.Now().Format(time.RFC3339)

		if err := helpers.SaveEventToSQLite(database, event); err != nil {
			log.Println("Failed to save event to SQLite:", err)
			return
		}

		newData, _ := json.Marshal(event)
		_, err = js.Publish("operation.stage1.data", newData)
		if err != nil {
			log.Println("Failed to publish to next stage:", err)
		}
	}, nats.Durable("patient-intake-durable"), nats.ManualAck())
	if err != nil {
		return fmt.Errorf("subscription error: %w", err)
	}
	defer func(sub *nats.Subscription) {
		_ = sub.Unsubscribe()
	}(sub)

	log.Println("[patient-intake] Listening on operation.incoming.data")
	<-cmd.Context().Done()

	return nil
}
