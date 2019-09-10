package songkick

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Billing string

const (
	BillingHeadline Billing = "headline"
	BillingSupport  Billing = "support"
)

type MusicBrainz struct {
	ID   string
	Href string
}

type Artist struct {
	DisplayName string        `json:"displayName"`
	ID          int           `json:"id"`
	URI         string        `json:"uri"`
	MusicBrainz []MusicBrainz `json:"indentifier"`
}

type City struct {
	DisplayName string `json:"displayName"`
}

type Venue struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	URI         string `json:"uri"`
	City        City   `json:"city"`
}

type venueResults struct {
	ResultsPage struct {
		Results struct {
			Venue Venue `json:"venue"`
		} `json:"results"`
	} `json:"resultsPage"`
}

type Event struct {
	Venue       Venue `json:"venue"`
	Performance []struct {
		Artist       Artist  `json:"artist"`
		Billing      Billing `json:"billing"`
		BillingIndex int     `json:"billingIndex"`
		DisplayName  string  `json:"displayName"`
	} `json:"performance"`
}

type eventsResults struct {
	ResultsPage struct {
		Page         int `json:"page"`
		TotalEntries int `json:"totalEntries"`
		Results      struct {
			Events []Event `json:"event"`
		} `json:"results"`
	} `json:"resultsPage"`
}

type GetVenueCalendarReq struct {
	VenueID string
	Page    int // Defaults to 1
	PerPage int // Default to 50
	// MinDate time.Time
	// MaxDate time.Time
}

type EventsResp struct {
	Events []Event `json:"events"`
}

type GetVenueReq struct {
	VenueID string
}

type Client interface {
	GetVenueCalendar(context.Context, *GetVenueCalendarReq) (*EventsResp, error)
	GetVenue(context.Context, *GetVenueReq) (*Venue, error)
}

type APIClient struct {
	Client http.Client
	Key    string
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

func NewClient(key string) *APIClient {
	return &APIClient{
		Key: key,
	}
}

var _ Client = (*APIClient)(nil)
