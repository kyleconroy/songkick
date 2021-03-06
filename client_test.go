package songkick

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestParsing(t *testing.T) {
	blob, err := ioutil.ReadFile("testdata/calendar.json")
	if err != nil {
		t.Fatal(err)
	}
	var resp eventsResults
	if err := json.Unmarshal(blob, &resp); err != nil {
		t.Fatal(err)
	}
	if resp.ResultsPage.TotalEntries != len(resp.ResultsPage.Results.Events) {
		t.Errorf("incorrect number of events")
	}
}

func TestClient(t *testing.T) {
	key := os.Getenv("SONGKICK_API_KEY")
	if key == "" {
		t.Skip("no api key")
	}
	client := NewClient(key)

	t.Run("get-venue-calendar", func(t *testing.T) {
		resp, err := client.GetVenueCalendar(context.Background(), &GetVenueCalendarReq{
			VenueID: "6239-fillmore",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Events) == 0 {
			t.Errorf("no events found for the fillmore? doubtful")
		}
	})

	t.Run("get-venue", func(t *testing.T) {
		id := "6239-fillmore"
		resp, err := client.GetVenue(context.Background(), &GetVenueReq{
			VenueID: id,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Log(resp)
	})

	t.Run("find-venues", func(t *testing.T) {
		resp, err := client.FindVenues(context.Background(), &FindVenuesReq{
			Query: "Paper Tiger",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(resp.Venues) == 0 {
			t.Errorf("no venues found for the Fillmore? doubtful")
		}
	})
}
