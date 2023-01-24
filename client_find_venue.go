package songkick

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type venuesResults struct {
	ResultsPage struct {
		Page         int    `json:"page"`
		PerPage      int    `json:"per_page"`
		TotalEntries int    `json:"totalEntries"`
		Status       string `json:"status"`
		Error        struct {
			Message string `json:"message"`
		} `json:"error"`
		Results struct {
			Venues []Venue `json:"venue"`
		} `json:"results"`
	} `json:"resultsPage"`
}

type VenuesResp struct {
	Venues []Venue `json:"venues"`
}

type FindVenuesReq struct {
	Query   string
	Page    int // Defaults to 1
	PerPage int // Default to 50
	// MinDate time.Time
	// MaxDate time.Time
}

func (c *APIClient) FindVenues(ctx context.Context, req *FindVenuesReq) (*VenuesResp, error) {
	if req == nil {
		return nil, fmt.Errorf("nil GetVenueReq")
	}

	base := "https://api.songkick.com/api/3.0/search/venues.json"
	r, _ := http.NewRequest("GET", base, nil)
	r.Header.Set("User-Agent", "songkick/v1.0.0 (+https://github.com/kyleconroy/songkick)")
	r = r.WithContext(ctx)

	q := r.URL.Query()
	q.Add("apikey", c.Key)
	q.Add("query", req.Query)
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
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var vr venuesResults
	if err := dec.Decode(&vr); err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %s message: %s code: %d", vr.ResultsPage.Status, vr.ResultsPage.Error.Message, resp.StatusCode)
	}

	venues := []Venue{}
	if len(vr.ResultsPage.Results.Venues) > 0 {
		venues = vr.ResultsPage.Results.Venues
	}
	return &VenuesResp{Venues: venues}, nil
}
