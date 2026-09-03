package trafficshaper

import (
	"context"
	"os"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestPipe(t *testing.T) {
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

	resource := &Pipe{
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

	id, err := controller.AddPipe(ctx, resource)
	if err != nil {
		t.Fatalf("AddPipe() error = %v", err)
	}

	t.Cleanup(func() {
		if err := controller.DeletePipe(ctx, id); err != nil {
			t.Errorf("DeletePipe() error = %v", err)
		}
	})

	got, err := controller.GetPipe(ctx, id)
	if err != nil {
		t.Fatalf("GetPipe() error = %v", err)
	}

	if got.Number == "" {
		t.Error("GetPipe().Number is empty")
	}

	if got.Enabled != resource.Enabled {
		t.Errorf("GetPipe().Enabled = %q, want %q", got.Enabled, resource.Enabled)
	}

	if got.Bandwidth != resource.Bandwidth {
		t.Errorf("GetPipe().Bandwidth = %q, want %q", got.Bandwidth, resource.Bandwidth)
	}

	if got.BandwidthMetric != resource.BandwidthMetric {
		t.Errorf(
			"GetPipe().BandwidthMetric = %q, want %q",
			got.BandwidthMetric,
			resource.BandwidthMetric,
		)
	}

	if got.Queue != resource.Queue {
		t.Errorf("GetPipe().Queue = %q, want %q", got.Queue, resource.Queue)
	}

	if got.Mask != resource.Mask {
		t.Errorf("GetPipe().Mask = %q, want %q", got.Mask, resource.Mask)
	}

	if got.Buckets != resource.Buckets {
		t.Errorf("GetPipe().Buckets = %q, want %q", got.Buckets, resource.Buckets)
	}

	if got.Scheduler != resource.Scheduler {
		t.Errorf("GetPipe().Scheduler = %q, want %q", got.Scheduler, resource.Scheduler)
	}

	if got.CoDelEnable != resource.CoDelEnable {
		t.Errorf("GetPipe().CoDelEnable = %q, want %q", got.CoDelEnable, resource.CoDelEnable)
	}

	if got.CoDelTarget != resource.CoDelTarget {
		t.Errorf("GetPipe().CoDelTarget = %q, want %q", got.CoDelTarget, resource.CoDelTarget)
	}

	if got.CoDelInterval != resource.CoDelInterval {
		t.Errorf("GetPipe().CoDelInterval = %q, want %q", got.CoDelInterval, resource.CoDelInterval)
	}

	if got.CoDelECNEnable != resource.CoDelECNEnable {
		t.Errorf(
			"GetPipe().CoDelECNEnable = %q, want %q",
			got.CoDelECNEnable,
			resource.CoDelECNEnable,
		)
	}

	if got.PIEEnable != resource.PIEEnable {
		t.Errorf("GetPipe().PIEEnable = %q, want %q", got.PIEEnable, resource.PIEEnable)
	}

	if got.FQCoDelQuantum != resource.FQCoDelQuantum {
		t.Errorf(
			"GetPipe().FQCoDelQuantum = %q, want %q",
			got.FQCoDelQuantum,
			resource.FQCoDelQuantum,
		)
	}

	if got.FQCoDelLimit != resource.FQCoDelLimit {
		t.Errorf(
			"GetPipe().FQCoDelLimit = %q, want %q",
			got.FQCoDelLimit,
			resource.FQCoDelLimit,
		)
	}

	if got.FQCoDelFlows != resource.FQCoDelFlows {
		t.Errorf(
			"GetPipe().FQCoDelFlows = %q, want %q",
			got.FQCoDelFlows,
			resource.FQCoDelFlows,
		)
	}

	if got.Origin == "" {
		t.Error("GetPipe().Origin is empty")
	}

	if got.Delay != resource.Delay {
		t.Errorf("GetPipe().Delay = %q, want %q", got.Delay, resource.Delay)
	}

	if got.Description != resource.Description {
		t.Errorf("GetPipe().Description = %q, want %q", got.Description, resource.Description)
	}

	resource.Description = "updated test pipe"

	if err := controller.UpdatePipe(ctx, id, resource); err != nil {
		t.Fatalf("UpdatePipe() error = %v", err)
	}

	got, err = controller.GetPipe(ctx, id)
	if err != nil {
		t.Fatalf("GetPipe() after UpdatePipe() error = %v", err)
	}

	if got.Description != resource.Description {
		t.Errorf(
			"GetPipe().Description after update = %q, want %q",
			got.Description,
			resource.Description,
		)
	}
}
