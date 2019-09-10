package songkick

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type venueResults struct {
	ResultsPage struct {
		Status  string `json:"status"`
		Results struct {
			Venue Venue `json:"venue"`
		} `json:"results"`
	} `json:"resultsPage"`
}

type GetVenueReq struct {
	VenueID string
}

func (c *APIClient) GetVenue(ctx context.Context, req *GetVenueReq) (*Venue, error) {
	if req == nil {
		return nil, fmt.Errorf("nil GetVenueReq")
	}

	base := "https://api.songkick.com/api/3.0/venues/%s.json"
	r, _ := http.NewRequest("GET", fmt.Sprintf(base, req.VenueID), nil)
	r.Header.Set("User-Agent", "songkick/v1.0.0 (+https://github.com/kyleconroy/songkick)")
	r = r.WithContext(ctx)

	q := r.URL.Query()
	q.Add("apikey", c.Key)
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
	var vr venueResults
	if err := dec.Decode(&vr); err != nil {
		return nil, err
	}
	return &vr.ResultsPage.Results.Venue, nil
}
