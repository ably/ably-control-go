package control

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountStats(t *testing.T) {
	client, _ := newTestClient(t)

	stats, err := client.AccountStats(nil)
	assert.NoError(t, err)
	assert.NotNil(t, stats)
}

func TestAccountStatsWithParams(t *testing.T) {
	client, _ := newTestClient(t)

	params := StatsQueryParams{
		Unit:      "hour",
		Direction: "backwards",
	}

	stats, err := client.AccountStats(&params)
	assert.NoError(t, err)
	assert.NotNil(t, stats)
}

func TestAppStats(t *testing.T) {
	client, _ := newTestClient(t)
	app := newTestApp(t, &client)

	stats, err := client.AppStats(app.ID, nil)
	assert.NoError(t, err)
	assert.NotNil(t, stats)
}

func TestAppStatsWithParams(t *testing.T) {
	client, _ := newTestClient(t)
	app := newTestApp(t, &client)

	params := StatsQueryParams{
		Unit:      "hour",
		Direction: "backwards",
	}

	stats, err := client.AppStats(app.ID, &params)
	assert.NoError(t, err)
	assert.NotNil(t, stats)
}

func TestStatsQueryParamsEncode(t *testing.T) {
	start := 1000
	end := 2000
	limit := 10

	params := StatsQueryParams{
		Start:     &start,
		End:       &end,
		Unit:      "hour",
		Direction: "backwards",
		Limit:     &limit,
	}

	encoded := params.encode()
	assert.Contains(t, encoded, "start=1000")
	assert.Contains(t, encoded, "end=2000")
	assert.Contains(t, encoded, "unit=hour")
	assert.Contains(t, encoded, "direction=backwards")
	assert.Contains(t, encoded, "limit=10")
	assert.True(t, encoded[0] == '?')
}

func TestStatsQueryParamsEncodeEmpty(t *testing.T) {
	params := StatsQueryParams{}
	assert.Equal(t, "", params.encode())
}

func TestStatsQueryParamsEncodePartial(t *testing.T) {
	params := StatsQueryParams{
		Unit: "day",
	}
	encoded := params.encode()
	assert.Equal(t, "?unit=day", encoded)
}
