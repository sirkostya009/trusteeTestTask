package service

import (
	"google.golang.org/api/calendar/v3"
	"slices"
	"testing"
	"time"
)

var service = &CinemaService{}

func TestFilterShowtimeEmpty(t *testing.T) {
	service.showtime = []*Showtime{newShowtime(1, time.Now().Truncate(time.Second), time.Hour)}
	var events []*calendar.Event

	filtered := service.FilterShowtime(events)

	if !slices.Equal(service.showtime, filtered) {
		t.Errorf("expected %v, got %v", service.showtime, filtered)
	}
}

func TestFilterShowtime(t *testing.T) {
	now := time.Now().Truncate(time.Second)

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
	now := time.Now().Truncate(time.Second)

	events := []*calendar.Event{
		newEvent(now.Add(time.Hour*2), now.Add(time.Hour*2+time.Minute*23)),
		newEvent(now, now.Add(time.Minute*23)),
	}

	if !service.IsShowtimeOverlapping(events, newShowtime(-1, now.Add(time.Hour*2), time.Hour)) {
		t.Error("expected true")
	}

	if service.IsShowtimeOverlapping(events, newShowtime(-1, now.Add(time.Hour), time.Hour)) {
		t.Error("overlapping should be inclusive")
	}
}

func newShowtime(id int, start time.Time, hours time.Duration) *Showtime {
	return &Showtime{
		Id:        id,
		BeginTime: start,
		EndTime:   start.Add(hours),
		Duration:  int(hours) * 60,
	}
}

func newEvent(start, end time.Time) *calendar.Event {
	return &calendar.Event{
		Start: &calendar.EventDateTime{DateTime: start.Format(time.RFC3339)},
		End:   &calendar.EventDateTime{DateTime: end.Format(time.RFC3339)},
	}
}
