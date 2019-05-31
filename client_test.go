package songkick

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestJSONParsing(t *testing.T) {
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

func TestRequest(t *testing.T) {
	key := os.Getenv("SONGKICK_API_KEY")
	if key == "" {
		t.Skip("no api key")
	}
	client := NewClient(key)
	resp, err := client.GetVenueCalendar(context.Background(), &GetVenueCalendarReq{
		VenueID: "6239-fillmore",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Events) == 0 {
		t.Errorf("no events found for the fillmore? doubtful")
	}
}
