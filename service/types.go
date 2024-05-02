package service

import "time"

type Showtime struct {
	Id        int               `json:"id"`
	BeginTime time.Time         `json:"beginTime"`
	EndTime   time.Time         `json:"endTime"`
	Duration  int               `json:"duration"` // in minutes
	Title     string            `json:"title"`
	Seats     int               `json:"seats"`
	ImgSrc    string            `json:"image"`
	bookers   map[string]string // ip -> eventId
}

func (s *Showtime) ToDto(requester string) ShowtimeDto {
	_, booked := s.bookers[requester]
	return ShowtimeDto{
		Id:        s.Id,
		BeginTime: s.BeginTime,
		EndTime:   s.EndTime,
		Duration:  s.Duration,
		Title:     s.Title,
		Seats:     s.Seats,
		IsBooked:  booked,
	}
}

type ShowtimeDto struct {
	Id        int       `json:"id"`
	BeginTime time.Time `json:"beginTime"`
	EndTime   time.Time `json:"endTime"`
	Duration  int       `json:"duration"` // in minutes
	Title     string    `json:"title"`
	Seats     int       `json:"seats"`
	IsBooked  bool      `json:"isBooked"`
}
