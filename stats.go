package control

import (
	"fmt"
	neturl "net/url"
	"strconv"
)

// StatsQueryParams contains query parameters for stats endpoints.
type StatsQueryParams struct {
	// Start of the query interval as a Unix timestamp (milliseconds).
	Start *int
	// End of the query interval as a Unix timestamp (milliseconds).
	End *int
	// The unit of time for stats aggregation (minute, hour, day, month).
	Unit string
	// The direction of the query (forwards or backwards).
	Direction string
	// The maximum number of records to return.
	Limit *int
}

func (p *StatsQueryParams) encode() string {
	v := neturl.Values{}
	if p.Start != nil {
		v.Set("start", strconv.Itoa(*p.Start))
	}
	if p.End != nil {
		v.Set("end", strconv.Itoa(*p.End))
	}
	if p.Unit != "" {
		v.Set("unit", p.Unit)
	}
	if p.Direction != "" {
		v.Set("direction", p.Direction)
	}
	if p.Limit != nil {
		v.Set("limit", strconv.Itoa(*p.Limit))
	}
	encoded := v.Encode()
	if encoded != "" {
		return "?" + encoded
	}
	return ""
}

// AccountStatsResponse represents a stats response for an account.
type AccountStatsResponse struct {
	// The interval ID for this stats entry.
	IntervalId string `json:"intervalId"`
	// The unit of time for this stats entry.
	Unit string `json:"unit"`
	// The schema version for the stats data.
	Schema string `json:"schema"`
	// The stats entries.
	Entries map[string]any `json:"entries"`
	// Whether these stats are still being calculated.
	InProgress string `json:"inProgress,omitempty"`
	// The account ID these stats belong to.
	AccountId string `json:"accountId"`
}

// AppStatsResponse represents a stats response for an app.
type AppStatsResponse struct {
	// The interval ID for this stats entry.
	IntervalId string `json:"intervalId"`
	// The unit of time for this stats entry.
	Unit string `json:"unit"`
	// The schema version for the stats data.
	Schema string `json:"schema"`
	// The stats entries.
	Entries map[string]any `json:"entries"`
	// Whether these stats are still being calculated.
	InProgress string `json:"inProgress,omitempty"`
	// The app ID these stats belong to.
	AppId string `json:"appId"`
}

// AccountStats retrieves stats for the account.
func (c *Client) AccountStats(params *StatsQueryParams) ([]AccountStatsResponse, error) {
	var stats []AccountStatsResponse
	path := fmt.Sprintf("/accounts/%s/stats", c.accountID)
	if params != nil {
		path += params.encode()
	}
	err := c.request("GET", path, nil, &stats)
	return stats, err
}

// AppStats retrieves stats for the specified application.
func (c *Client) AppStats(appID string, params *StatsQueryParams) ([]AppStatsResponse, error) {
	var stats []AppStatsResponse
	path := fmt.Sprintf("/apps/%s/stats", appID)
	if params != nil {
		path += params.encode()
	}
	err := c.request("GET", path, nil, &stats)
	return stats, err
}
