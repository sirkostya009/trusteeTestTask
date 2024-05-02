package service

import (
	"errors"
	"google.golang.org/api/calendar/v3"
	"math/rand/v2"
	"slices"
	"strconv"
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

	now := time.Now().Truncate(time.Second)

	s := make([]*Showtime, 10)
	for i := range s {
		hour := time.Duration(rand.Int64N(14) + 8)
		day := time.Duration(rand.Int64N(3))
		hours := rand.IntN(3) + 1
		titleI := rand.IntN(len(movieTitles))
		s[i] = &Showtime{
			Id:        i + 1,
			BeginTime: now.Add(time.Hour * hour * day),
			EndTime:   now.Add(time.Hour * hour * day).Add(time.Hour * time.Duration(hours)),
			Duration:  hours * 60,
			Title:     movieTitles[titleI],
			Seats:     rand.IntN(100) + 1,
			ImgSrc:    "image" + strconv.Itoa(i) + ".jpg",
			bookers:   map[string]string{},
		}
	}
	slices.SortFunc(s, func(a, b *Showtime) int {
		if cmp := a.BeginTime.Compare(b.BeginTime); cmp != 0 {
			return cmp
		}
		return a.EndTime.Compare(b.EndTime)
	})
	last := s[len(s)-1].EndTime
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

func (s *CinemaService) IsShowtimeOverlapping(events []*calendar.Event, st *Showtime) bool {
	for _, e := range events {
		evStartTime, err := time.Parse(time.RFC3339, e.Start.DateTime)
		if err != nil {
			continue
		}
		evEndTime, err := time.Parse(time.RFC3339, e.End.DateTime)
		if err != nil {
			continue
		}

		// uncomment for debugging purposes
		//fmt.Printf("%v | %v | %v | %v\n%v | %v | %v\n",
		//	st.BeginTime.Format(time.TimeOnly),
		//	st.EndTime.Format(time.TimeOnly),
		//	evStartTime.Format(time.TimeOnly),
		//	evEndTime.Format(time.TimeOnly),
		//	st.BeginTime.Compare(evStartTime) >= 0 && st.BeginTime.Before(evEndTime),
		//	st.EndTime.After(evStartTime) && st.EndTime.Before(evEndTime),
		//	st.BeginTime.Compare(evStartTime) <= 0 && st.EndTime.Compare(evEndTime) >= 0,
		//)

		if st.BeginTime.Compare(evStartTime) >= 0 && st.BeginTime.Before(evEndTime) || // starts during
			st.EndTime.After(evStartTime) && st.EndTime.Before(evEndTime) || // ends during
			st.BeginTime.Compare(evStartTime) <= 0 && st.EndTime.Compare(evEndTime) >= 0 /*starts before and ends after*/ {
			return true
		}
	}
	return false
}

func (s *CinemaService) GetShowtime(id int) (*Showtime, error) {
	for _, st := range s.showtime {
		if st.Id == id {
			return st, nil
		}
	}
	return nil, errors.New("showtime not found")
}

func (s *CinemaService) BookShowtime(id int, ip, eventId string) error {
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

func (s *CinemaService) CancelBooking(id int, ip string) string {
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

func (s *CinemaService) IsBooked(id int, ip string) bool {
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
