package adaptersapi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testit-tms/adapters-go/v2/adaptersapi"
)

func TestAutoTestCreateApiModel_layer_omitted_when_empty(t *testing.T) {
	req := adaptersapi.NewAutoTestCreateApiModel("project", "ext", "name")
	data, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotContains(t, string(data), `"layer"`)
}

func TestAutoTestCreateApiModel_layer_with_source_run(t *testing.T) {
	req := adaptersapi.NewAutoTestCreateApiModel("project", "ext", "name")
	req.SetLayer(*adaptersapi.NewLayerApiModel("API", adaptersapi.LAYERSOURCE_RUN))

	data, err := json.Marshal(req)
	require.NoError(t, err)
	require.Contains(t, string(data), `"layer":{"name":"API","source":"Run"}`)
}

func TestAutoTestUpdateApiModel_resetLayer_false_and_layer(t *testing.T) {
	req := adaptersapi.NewAutoTestUpdateApiModel("project", "ext", "name", false)
	req.SetLayer(*adaptersapi.NewLayerApiModel("my-custom-layer", adaptersapi.LAYERSOURCE_RUN))

	data, err := json.Marshal(req)
	require.NoError(t, err)
	require.Contains(t, string(data), `"resetLayer":false`)
	require.Contains(t, string(data), `"layer":{"name":"my-custom-layer","source":"Run"}`)
}

func TestAutoTestUpdateApiModel_no_layer_resetLayer_false(t *testing.T) {
	req := adaptersapi.NewAutoTestUpdateApiModel("project", "ext", "name", false)

	data, err := json.Marshal(req)
	require.NoError(t, err)
	require.Contains(t, string(data), `"resetLayer":false`)
	require.NotContains(t, string(data), `"layer"`)
}
