package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dyammarcano/alfanumeric-cnpj/pkg/cnpj"
	"github.com/inovacc/ksuid"
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
		return fmt.Errorf("database connection error: %w", err)
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
		Subjects: []string{"operation.request.data", "operation.incoming.data"},
		Storage:  nats.FileStorage,
	})
	if err != nil && !strings.Contains(err.Error(), "stream name already in use") {
		return fmt.Errorf("stream creation error: %w", err)
	}

	// Create fake patient event
	eventData := common.PatientEvent{
		PatientID: ksuid.NewString(),
		CurrentID: ksuid.NewString(), // mutate in every step
		Step:      "patient-intake",  // mutate in every step
		Status:    "initiated",       // mutate in every step
		Timestamp: time.Now().Format(time.RFC3339),
		Patient: common.Patient{
			ID: cnpj.FormatCNPJ(cnpj.GenerateCNPJ()),
		},
	}

	msg := &nats.Msg{
		Subject: common.SubjectOperationRequestData,
		Data:    eventData.Bytes(),
		Header:  nats.Header{},
	}

	msg.Header.Add("current_id", eventData.CurrentID)
	msg.Header.Add("patient_id", eventData.PatientID)

	// Send as request to operation.request.data and wait reply
	replyMsg, err := natsConn.RequestMsgWithContext(cmd.Context(), msg)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return fmt.Errorf("no responders: %w", err)
		}
		return fmt.Errorf("request error: %w", err)
	}

	// Use response (if needed)
	log.Println("Received reply from request.data:", string(replyMsg.Data))

	// Now continue flow to operation.incoming.data
	_, err = js.Publish(common.SubjectOperationIncomingData, []byte(eventData.String()))
	if err != nil {
		return fmt.Errorf("publish error: %w", err)
	}

	log.Printf("Event sent to %s", common.SubjectOperationIncomingData)
	return nil
}
