package function

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Handle an HTTP Request.
/* func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain")

	_, err := fmt.Fprintf(res, "OK\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error or response write: %v", err)
	}
} */

type EventOrchestrator struct {
	SvcUrl string `json:"service_url"`
}

func Handle(ctx context.Context, event cloudevents.Event) (resp *cloudevents.Event, err error) {
	evOrch := &EventOrchestrator{}
	if err = event.DataAs(evOrch); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse incoming CloudEvent %s\n", err)
		return nil, err
	}

	// this is needed for creating a valid CloudEvent
	response := cloudevents.NewEvent()
	response.SetID("example-uuid-32943bac6fea")
	response.SetSource("purchase/getter")
	response.SetType("purchase")

	svcResp, err := http.Get(evOrch.SvcUrl)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(svcResp.Body)
	if err != nil {
		return nil, err
	}

	// Set the data from EventOrchestrator type
	response.SetData(cloudevents.ApplicationJSON, body)

	// Validate the response
	resp = &response
	if err = resp.Validate(); err != nil {
		fmt.Printf("invalid event created. %v", err)
	}

	return
}
