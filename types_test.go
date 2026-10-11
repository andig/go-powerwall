package powerwall

import (
	"encoding/json"
	"testing"
	"time"
)

func TestStatusDataNullTimes(t *testing.T) {
	// Some gateway emulators (e.g. pypowerwall in TEDAPI or cloud mode) report
	// start_time and up_time_seconds as null when they are unknown.
	var s StatusData
	err := json.Unmarshal([]byte(`{"din":"X","start_time":null,"up_time_seconds":null,"version":"26.26.11 c4db63ce"}`), &s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.StartTime.IsZero() || s.UpTime.Duration != 0 {
		t.Errorf("null should leave the values unset, got %v and %v", s.StartTime, s.UpTime)
	}
	if s.Version != "26.26.11 c4db63ce" {
		t.Errorf("other fields must still decode, got version %q", s.Version)
	}
}

func TestStatusDataGatewayFormats(t *testing.T) {
	var s StatusData
	err := json.Unmarshal([]byte(`{"start_time":"2024-03-11 09:12:41 +0800","up_time_seconds":"127h34m16.275122187s"}`), &s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := time.Date(2024, 3, 11, 9, 12, 41, 0, time.FixedZone("", 8*3600))
	if !s.StartTime.Equal(want) {
		t.Errorf("start_time = %v, want %v", s.StartTime, want)
	}
	if s.UpTime.Duration != 127*time.Hour+34*time.Minute+16275122187*time.Nanosecond {
		t.Errorf("up_time_seconds = %v", s.UpTime)
	}
}

func TestInvalidTimesStillFail(t *testing.T) {
	var s StatusData
	if err := json.Unmarshal([]byte(`{"start_time":"2026-07-17T15:51:37-07:00"}`), &s); err == nil {
		t.Error("an ISO-8601 start_time should still be rejected")
	}
	if err := json.Unmarshal([]byte(`{"up_time_seconds":"soon"}`), &s); err == nil {
		t.Error("an invalid duration should still be rejected")
	}
}
