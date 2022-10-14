package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/appuio/appuio-cloud-reporting/pkg/report"
	"github.com/appuio/appuio-cloud-reporting/pkg/sourcekey"
	"github.com/getlago/lago-go-client"
	"github.com/prometheus/common/model"

	"github.com/bastjan/lago-send-metrics-poc/lagoerr"
)

// Used to generate the unique transaction ID
const transactionNamespace = "appuio-cloud-reporting"

// Used to generate the unique transaction ID
// Increase to reset transaction IDs
const transactionEpoch = "11"

var q report.PromQuerier = promData
var apiKey = os.Getenv("LAGO_API_KEY")

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	lagoClient := lago.New().
		SetApiKey(apiKey).
		// SetDebug will log the RAW request and RAW response
		SetDebug(true)

	queryTime := time.Now().Truncate(time.Hour)
	res, warnings, err := q.Query(ctx, "PRODUCTS", queryTime)
	if err != nil {
		panic(err)
	}
	if len(warnings) > 0 {
		fmt.Println("Warnings:", warnings)
	}
	samples := res.(model.Vector)

	for _, sample := range samples {
		if err := processSample(ctx, lagoClient, sample); err != nil {
			panic(err)
		}
	}
}

func processSample(ctx context.Context, lagoClient *lago.Client, sample *model.Sample) error {
	fmt.Println("Inserting metric raw", sample.Metric, sample.Value)

	product, _ := sourcekey.Parse(string(sample.Metric["product"]))

	code := product.Query + ":" + product.Zone + ":" + product.Namespace + ":" + product.Class

	ev := lago.EventInput{
		TransactionID:      transactionID(sample.Timestamp.Time(), string(sample.Metric["product"])),
		ExternalCustomerID: string(sample.Metric["tenant"]),
		Timestamp:          time.Now().Unix(),
		Code:               code,
		Properties: map[string]string{
			"grams":  fmt.Sprintf("%f", sample.Value),
			"mymeta": "myvalue",
		},
	}
	fmt.Println("Inserting metered event", ev)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return storeEvent(ctx, lagoClient, ev)
}

func storeEvent(ctx context.Context, lagoClient *lago.Client, ev lago.EventInput) error {
	err := lagoClient.Event().Create(ctx, &ev)
	if err != nil {
		return fmt.Errorf("error creating event: %w", lagoerr.Wrap(err))
	}

	// Check if the event was created successfully
	// The API does not return an error on conflict or if an event has invalid data.
	// The API is also eventually consistent, so we need to wait a bit for the event to be created.
	for {
		storedEvent, err := lagoClient.Event().Get(ctx, ev.TransactionID)
		if err != nil {
			if err.HTTPStatusCode == 404 {
				if ctx.Err() != nil {
					return fmt.Errorf("timeout waiting for event to be stored: %w", lagoerr.Wrap(err))
				}
				continue
			}
			return lagoerr.Wrap(err)
		}
		fmt.Println("Stored event: ", storedEvent)
		break
	}
	return nil
}

func transactionID(t time.Time, product string) string {
	i := strings.Join([]string{transactionNamespace, transactionEpoch, t.UTC().Format(time.RFC3339), product}, ":")
	return fmt.Sprintf("%x", sha256.Sum256([]byte(i)))
}
