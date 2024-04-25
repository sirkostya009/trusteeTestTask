package service

import (
	"google.golang.org/api/calendar/v3"
	"slices"
	"testing"
	"time"
)

var service = &CinemaService{}

func TestFilterShowtimeEmpty(t *testing.T) {
	service.showtime = []*Showtime{newShowtime(1, time.Now(), time.Hour)}
	var events []*calendar.Event

	filtered := service.FilterShowtime(events)

	if !slices.Equal(service.showtime, filtered) {
		t.Errorf("expected %v, got %v", service.showtime, filtered)
	}
}

func TestFilterShowtime(t *testing.T) {
	now := time.Now()

	service.showtime = []*Showtime{
		newShowtime(2, now.Add(time.Hour*2), time.Hour),
		newShowtime(1, now, time.Hour),
		newShowtime(4, now.Add(time.Hour*8), time.Hour),
		newShowtime(3, now.Add(time.Hour*5), time.Hour),
	}
	events := []*calendar.Event{
		newEvent(now.Add(time.Hour*2), now.Add(time.Hour*2+time.Minute*23)),
		newEvent(now, now.Add(time.Minute*23)),
	}

	filtered := service.FilterShowtime(events)

	if len(filtered) != 2 || filtered[0].Id != 4 || filtered[1].Id != 3 {
		t.Errorf("expected %v, got %v", service.showtime[2:4], filtered)
	}
}

func TestIsShowtimeOverlapping(t *testing.T) {
	now := time.Now()

	events := []*calendar.Event{
		newEvent(now.Add(time.Hour*2), now.Add(time.Hour*2+time.Minute*23)),
		newEvent(now, now.Add(time.Minute*23)),
	}

	if !service.IsShowtimeOverlapping(events, newShowtime(-1, now.Add(time.Hour*2), time.Hour)) {
		t.Error("expected true")
	}

	// this fails because back-to-back timing is considered overlapping
	// due to non-inclusive nature of comparison by Before and After methods
	if service.IsShowtimeOverlapping(events, newShowtime(-1, now.Add(time.Hour), time.Hour)) {
		t.Error("expected false")
	}
}

func newShowtime(id int64, start time.Time, duration time.Duration) *Showtime {
	return &Showtime{Id: id, BeginTime: start, Duration: duration}
}

func newEvent(start, end time.Time) *calendar.Event {
	return &calendar.Event{
		Start: &calendar.EventDateTime{DateTime: start.Format(time.RFC3339)},
		End:   &calendar.EventDateTime{DateTime: end.Format(time.RFC3339)},
	}
}
