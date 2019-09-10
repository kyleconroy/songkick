package songkick

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type eventsResults struct {
	ResultsPage struct {
		Page         int    `json:"page"`
		PerPage      int    `json:"per_page"`
		TotalEntries int    `json:"totalEntries"`
		Status       string `json:"status"`
		Results      struct {
			Events []Event `json:"event"`
		} `json:"results"`
	} `json:"resultsPage"`
}

type EventsResp struct {
	Events []Event `json:"events"`
}

func (c *APIClient) GetVenueCalendar(ctx context.Context, req *GetVenueCalendarReq) (*EventsResp, error) {
	if req == nil {
		return nil, fmt.Errorf("nil GetVenueCalendarReq")
	}

	base := "https://api.songkick.com/api/3.0/venues/%s/calendar.json"
	r, _ := http.NewRequest("GET", fmt.Sprintf(base, req.VenueID), nil)
	r.Header.Set("User-Agent", "songkick/v1.0.0 (+https://github.com/kyleconroy/songkick)")
	r = r.WithContext(ctx)

	q := r.URL.Query()
	q.Add("apikey", c.Key)
	if req.Page != 0 {
		q.Add("page", strconv.Itoa(req.Page))
	}
	if req.PerPage != 0 {
		q.Add("per_page", strconv.Itoa(req.PerPage))
	}
	r.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-status")
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var er eventsResults
	if err := dec.Decode(&er); err != nil {
		return nil, err
	}
	// Check status field
	events := []Event{}
	if len(er.ResultsPage.Results.Events) > 0 {
		events = er.ResultsPage.Results.Events
	}
	return &EventsResp{Events: events}, nil
}
