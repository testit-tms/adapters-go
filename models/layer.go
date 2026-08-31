package models

import "strings"

// LayerSourceRun is the layer source used by adapters when sending layer from a test run.
const LayerSourceRun = "Run"

// TestLayers lists recommended autotest pyramid layer names (not enforced).
var TestLayers = struct {
	E2E         string
	UI          string
	API         string
	Contract    string
	Integration string
	Component   string
	Unit        string
}{
	E2E:         "E2E",
	UI:          "UI",
	API:         "API",
	Contract:    "Contract",
	Integration: "Integration",
	Component:   "Component",
	Unit:        "Unit",
}

// NormalizeLayer returns trimmed layer name and whether it should be sent to API.
func NormalizeLayer(layer string) (string, bool) {
	layer = strings.TrimSpace(layer)
	if layer == "" {
		return "", false
	}
	return layer, true
}
