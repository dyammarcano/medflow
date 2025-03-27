package service

import (
	"fmt"
	"github.com/dyammarcano/alfanumeric-cnpj/pkg/cnpj"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"log"
	"medflow/internal/common"
	"medflow/internal/helpers"
	"strings"
)

func StartOperationRequestHandler(cmd *cobra.Command, _ []string) error {
	natsConn, js, err := helpers.ConnectToNATS()
	if err != nil {
		return fmt.Errorf("nats connection error: %w", err)
	}
	defer natsConn.Close()

	// Ensure stream exists for request subject
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "OPERATION_REQUEST_STREAM",
		Subjects: []string{common.SubjectOperationRequestData},
		Storage:  nats.FileStorage,
	})
	if err != nil && !strings.Contains(err.Error(), "subjects overlap with an existing stream") {
		return fmt.Errorf("stream creation error: %w", err)
	}

	_, err = natsConn.QueueSubscribe(common.SubjectOperationRequestData, "request-workers", func(msg *nats.Msg) {
		currentID := msg.Header.Get("current_id")
		patientID := msg.Header.Get("patient_id")

		if currentID == "" || patientID == "" {
			log.Println("Missing required headers")
			return
		}

		var eventData common.PatientEvent
		if err := eventData.Decode(msg.Data); err != nil {
			log.Println("Invalid request JSON:", err)
			return
		}

		log.Printf("[operation-request] Received request for: %s", eventData.PatientID)

		if !cnpj.IsValid(eventData.Patient.ID) {
			log.Println("Invalid patient ID")
			return
		}

		eventResponse := helpers.GeneratePatient()
		eventResponse.CurrentID = currentID
		eventResponse.Step = "operation-request"
		eventResponse.Status = "success"

		msg = &nats.Msg{
			Subject: common.SubjectOperationRequestData,
			Header:  nats.Header{},
		}

		msg.Header.Add("currentI_id", eventData.CurrentID)
		msg.Header.Add("patient_id", eventResponse.PatientID)

		// Echo back the same event as response
		if err := msg.RespondMsg(msg); err != nil {
			log.Println("Error sending response:", err)
		} else {
			log.Printf("[operation-request] Replied to request for: %s", eventData.PatientID)
		}
	})
	if err != nil {
		return fmt.Errorf("subscription error: %w", err)
	}

	log.Println("[operation-request] Listening for operation.request.data requests")
	<-cmd.Context().Done()
	return nil
}
