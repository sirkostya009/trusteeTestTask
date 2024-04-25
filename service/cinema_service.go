package service

import (
	"errors"
	"google.golang.org/api/calendar/v3"
	"math/rand/v2"
	"slices"
	"sync"
	"time"
)

type CinemaService struct {
	lock               sync.Mutex
	showtime           []*Showtime
	maxShowtimeEndTime time.Time
}

func NewCinemaService() *CinemaService {
	movieTitles := []string{
		"The Matrix",
		"The Shawshank Redemption",
		"Schindler's List",
		"Star Wars",
		"Straight Outta Compton",
		"Interstellar",
		"The Dark Knight",
		"Back To The Future",
		"20 Days in Mariupol",
		"Dune",
	}

	s := make([]*Showtime, 10)
	for i := range s {
		hour := time.Duration(rand.Int64N(14) + 8)
		day := time.Duration(rand.Int64N(3))
		hours := time.Duration(rand.Int64N(3) + 1)
		titleI := rand.IntN(len(movieTitles))
		s[i] = &Showtime{
			Id:        int64(i + 1),
			BeginTime: time.Now().Add(time.Hour * hour * day),
			Duration:  time.Hour * hours,
			Title:     movieTitles[titleI],
			Seats:     rand.IntN(100) + 1,
			bookers:   map[string]string{},
		}
	}
	slices.SortFunc(s, func(a, b *Showtime) int {
		return a.BeginTime.Compare(b.BeginTime)
	})
	last := s[len(s)-1].BeginTime.Add(s[len(s)-1].Duration)
	return &CinemaService{showtime: s, maxShowtimeEndTime: last}
}

func (s *CinemaService) MaxShowtimeEndTime() time.Time {
	return s.maxShowtimeEndTime
}

func (s *CinemaService) Showtime() []*Showtime {
	return s.showtime
}

func (s *CinemaService) FilterShowtime(events []*calendar.Event) []*Showtime {
	if len(events) == 0 {
		return s.showtime
	}

	var filtered []*Showtime

	for _, showtime := range s.showtime {
		if !s.IsShowtimeOverlapping(events, showtime) {
			filtered = append(filtered, showtime)
		}
	}

	return filtered
}

func (s *CinemaService) IsShowtimeOverlapping(events []*calendar.Event, showtime *Showtime) bool {
	for _, e := range events {
		startTime, err := time.Parse(time.RFC3339, e.Start.DateTime)
		if err != nil {
			continue
		}
		endTime, err := time.Parse(time.RFC3339, e.End.DateTime)
		if err != nil {
			continue
		}

		if (showtime.BeginTime.After(startTime) && showtime.BeginTime.Before(endTime)) ||
			(showtime.BeginTime.Before(startTime) && showtime.BeginTime.Add(showtime.Duration).After(startTime)) {
			return true
		}
	}
	return false
}

func (s *CinemaService) GetShowtime(id int64) (*Showtime, error) {
	for _, st := range s.showtime {
		if st.Id == id {
			return st, nil
		}
	}
	return nil, errors.New("showtime not found")
}

func (s *CinemaService) BookShowtime(id int64, ip, eventId string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, st := range s.showtime {
		if st.Id == id {
			st.Seats--
			st.bookers[ip] = eventId
			return nil
		}
	}
	return errors.New("showtime not found")
}

func (s *CinemaService) CancelBooking(id int64, ip string) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, st := range s.showtime {
		if st.Id == id {
			st.Seats++
			eventId := st.bookers[ip]
			delete(st.bookers, ip)
			return eventId
		}
	}
	return ""
}

func (s *CinemaService) IsBooked(id int64, ip string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, st := range s.showtime {
		if st.Id == id {
			_, ok := st.bookers[ip]
			return ok
		}
	}
	return false
}
