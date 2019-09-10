package songkick

import (
	"context"
	"net/http"
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

type Event struct {
	Venue       Venue `json:"venue"`
	Performance []struct {
		Artist       Artist  `json:"artist"`
		Billing      Billing `json:"billing"`
		BillingIndex int     `json:"billingIndex"`
		DisplayName  string  `json:"displayName"`
	} `json:"performance"`
}

type GetVenueCalendarReq struct {
	VenueID string
	Page    int // Defaults to 1
	PerPage int // Default to 50
	// MinDate time.Time
	// MaxDate time.Time
}

type Client interface {
	GetVenueCalendar(context.Context, *GetVenueCalendarReq) (*EventsResp, error)
	GetVenue(context.Context, *GetVenueReq) (*Venue, error)
	FindVenues(context.Context, *FindVenuesReq) (*VenuesResp, error)
}

type APIClient struct {
	Client http.Client
	Key    string
}

func NewClient(key string) *APIClient {
	return &APIClient{
		Key: key,
	}
}

var _ Client = (*APIClient)(nil)
