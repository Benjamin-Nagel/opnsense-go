package trafficshaper

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestQueue(t *testing.T) {
	opnsenseURL := os.Getenv("OPNSENSE_URI")
	opnsenseKey := os.Getenv("OPNSENSE_API_KEY")
	opnsenseSecret := os.Getenv("OPNSENSE_API_SECRET")

	apiClient := api.NewClient(api.Options{
		Uri:           opnsenseURL,
		APIKey:        opnsenseKey,
		APISecret:     opnsenseSecret,
		AllowInsecure: true,
		MaxBackoff:    30,
		MinBackoff:    1,
		MaxRetries:    4,
	})

	controller := Controller{Api: apiClient}
	ctx := context.Background()

	pipe := &Pipe{
		Number:          "1",
		Enabled:         "1",
		Bandwidth:       "100",
		BandwidthMetric: api.SelectedMap("Kbit"),
		Queue:           "2",
		Mask:            api.SelectedMap("none"),
		Buckets:         "16",
		Scheduler:       api.SelectedMap("fifo"),
		CoDelEnable:     "1",
		CoDelTarget:     "5",
		CoDelInterval:   "100",
		CoDelECNEnable:  "1",
		PIEEnable:       "0",
		FQCoDelQuantum:  "1514",
		FQCoDelLimit:    "10240",
		FQCoDelFlows:    "1024",
		Origin:          "user",
		Delay:           "1",
		Description:     "test pipe",
	}

	pipeId, err := controller.AddPipe(ctx, pipe)
	if err != nil {
		t.Fatalf("AddPipe() error = %v", err)
	}

	t.Cleanup(func() {
		if err := controller.DeletePipe(ctx, pipeId); err != nil {
			t.Errorf("DeletePipe() error = %v", err)
		}
	})

	resource := &Queue{
		Number:         "1",
		Enabled:        "1",
		Pipe:           api.SelectedMap(pipeId),
		Weight:         "10",
		Mask:           api.SelectedMap("none"),
		Buckets:        "16",
		CoDelEnable:    "1",
		CoDelTarget:    "5",
		CoDelInterval:  "100",
		CoDelECNEnable: "1",
		PIEEnable:      "0",
		Description:    "test queue",
		Origin:         "user",
	}

	id, err := controller.AddQueue(ctx, resource)
	if err != nil {
		t.Fatalf("AddQueue() error = %v", err)
	}

	t.Cleanup(func() {
		if err := controller.DeleteQueue(ctx, id); err != nil {
			t.Errorf("DeleteQueue() error = %v", err)
		}
	})

	got, err := controller.GetQueue(ctx, id)
	if err != nil {
		t.Fatalf("GetQueue() error = %v", err)
	}

	if got.Number == "" {
		t.Error("GetQueue().Number is empty")
	}

	if got.Enabled != resource.Enabled {
		t.Errorf("GetQueue().Enabled = %q, want %q", got.Enabled, resource.Enabled)
	}

	if got.Pipe != resource.Pipe {
		t.Errorf("GetQueue().Pipe = %q, want %q", got.Pipe, resource.Pipe)
	}

	if got.Weight != resource.Weight {
		t.Errorf("GetQueue().Weight = %q, want %q", got.Weight, resource.Weight)
	}

	if got.Mask != resource.Mask {
		t.Errorf("GetQueue().Mask = %q, want %q", got.Mask, resource.Mask)
	}

	if got.Buckets != resource.Buckets {
		t.Errorf("GetQueue().Buckets = %q, want %q", got.Buckets, resource.Buckets)
	}

	if got.CoDelEnable != resource.CoDelEnable {
		t.Errorf("GetQueue().CoDelEnable = %q, want %q", got.CoDelEnable, resource.CoDelEnable)
	}

	if got.CoDelTarget != resource.CoDelTarget {
		t.Errorf("GetQueue().CoDelTarget = %q, want %q", got.CoDelTarget, resource.CoDelTarget)
	}

	if got.CoDelInterval != resource.CoDelInterval {
		t.Errorf("GetQueue().CoDelInterval = %q, want %q", got.CoDelInterval, resource.CoDelInterval)
	}

	if got.CoDelECNEnable != resource.CoDelECNEnable {
		t.Errorf(
			"GetQueue().CoDelECNEnable = %q, want %q",
			got.CoDelECNEnable,
			resource.CoDelECNEnable,
		)
	}

	if got.PIEEnable != resource.PIEEnable {
		t.Errorf("GetQueue().PIEEnable = %q, want %q", got.PIEEnable, resource.PIEEnable)
	}

	if got.Description != resource.Description {
		t.Errorf("GetQueue().Description = %q, want %q", got.Description, resource.Description)
	}

	if got.Origin == "" {
		t.Error("GetQueue().Origin is empty")
	}

	resource.Description = "updated test queue"

	if err := controller.UpdateQueue(ctx, id, resource); err != nil {
		t.Fatalf("UpdateQueue() error = %v", err)
	}

	got, err = controller.GetQueue(ctx, id)
	if err != nil {
		t.Fatalf("GetQueue() after UpdateQueue() error = %v", err)
	}

	if got.Description != resource.Description {
		t.Errorf(
			"GetQueue().Description after update = %q, want %q",
			got.Description,
			resource.Description,
		)
	}
}
