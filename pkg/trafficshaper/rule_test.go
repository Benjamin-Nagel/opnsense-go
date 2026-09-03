package trafficshaper

import (
	"context"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/browningluke/opnsense-go/pkg/api"
)

func TestRule(t *testing.T) {
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

	resource := &Rule{
		Enabled:         "1",
		Sequence:        "10",
		Interface:       api.SelectedMap("wan"),
		Protocol:        api.SelectedMap("tcp"),
		IPLength:        "10000",
		Source:          api.SelectedMapList{"192.168.1.0/24"},
		SourceNot:       "0",
		SourcePort:      "443",
		Destination:     api.SelectedMapList{"10.0.0.0/24"},
		DestinationNot:  "0",
		DestinationPort: "443",
		DSCP:            api.SelectedMapList{"af32", "af42"},
		Direction:       api.SelectedMap("in"),
		Target:          api.SelectedMap(pipeId),
		Description:     "test rule",
		Origin:          "user",
	}

	id, err := controller.AddRule(ctx, resource)
	if err != nil {
		t.Fatalf("AddRule() error = %v", err)
	}

	t.Cleanup(func() {
		if err := controller.DeleteRule(ctx, id); err != nil {
			t.Errorf("DeleteRule() error = %v", err)
		}
	})

	got, err := controller.GetRule(ctx, id)
	if err != nil {
		t.Fatalf("GetRule() error = %v", err)
	}

	if got.Enabled != resource.Enabled {
		t.Errorf("GetRule().Enabled = %q, want %q", got.Enabled, resource.Enabled)
	}

	if got.Sequence != resource.Sequence {
		t.Errorf("GetRule().Sequence = %q, want %q", got.Sequence, resource.Sequence)
	}

	if got.Interface != resource.Interface {
		t.Errorf("GetRule().Interface = %q, want %q", got.Interface, resource.Interface)
	}

	// there is no second interface in the test environment - skip
	// if got.Interface2 != resource.Interface2 {
	// 	t.Errorf("GetRule().Interface2 = %q, want %q", got.Interface2, resource.Interface2)
	// }

	if got.Protocol != resource.Protocol {
		t.Errorf("GetRule().Protocol = %q, want %q", got.Protocol, resource.Protocol)
	}

	if got.IPLength != resource.IPLength {
		t.Errorf("GetRule().IPLength = %q, want %q", got.IPLength, resource.IPLength)
	}

	assertSelectedMapListEqual(t, "Source", got.Source, resource.Source)

	if got.SourceNot != resource.SourceNot {
		t.Errorf("GetRule().SourceNot = %q, want %q", got.SourceNot, resource.SourceNot)
	}

	if got.SourcePort != resource.SourcePort {
		t.Errorf("GetRule().SourcePort = %q, want %q", got.SourcePort, resource.SourcePort)
	}

	assertSelectedMapListEqual(t, "Destination", got.Destination, resource.Destination)

	if got.DestinationNot != resource.DestinationNot {
		t.Errorf(
			"GetRule().DestinationNot = %q, want %q",
			got.DestinationNot,
			resource.DestinationNot,
		)
	}

	if got.DestinationPort != resource.DestinationPort {
		t.Errorf(
			"GetRule().DestinationPort = %q, want %q",
			got.DestinationPort,
			resource.DestinationPort,
		)
	}

	assertSelectedMapListEqual(t, "DSCP", got.DSCP, resource.DSCP)

	if got.Direction != resource.Direction {
		t.Errorf("GetRule().Direction = %q, want %q", got.Direction, resource.Direction)
	}

	if got.Target != resource.Target {
		t.Errorf("GetRule().Target = %q, want %q", got.Target, resource.Target)
	}

	if got.Description != resource.Description {
		t.Errorf("GetRule().Description = %q, want %q", got.Description, resource.Description)
	}

	if got.Origin == "" {
		t.Error("GetRule().Origin is empty")
	}

	resource.Description = "updated test rule"

	if err := controller.UpdateRule(ctx, id, resource); err != nil {
		t.Fatalf("UpdateRule() error = %v", err)
	}

	got, err = controller.GetRule(ctx, id)
	if err != nil {
		t.Fatalf("GetRule() after UpdateRule() error = %v", err)
	}

	if got.Description != resource.Description {
		t.Errorf(
			"GetRule().Description after update = %q, want %q",
			got.Description,
			resource.Description,
		)
	}
}

func assertSelectedMapListEqual(t *testing.T, field string, got, want api.SelectedMapList) {
	t.Helper()

	gotSorted := append(api.SelectedMapList(nil), got...)
	wantSorted := append(api.SelectedMapList(nil), want...)

	sort.Strings(gotSorted)
	sort.Strings(wantSorted)

	if !reflect.DeepEqual(gotSorted, wantSorted) {
		t.Errorf("GetRule().%s = %v, want %v", field, got, want)
	}
}
