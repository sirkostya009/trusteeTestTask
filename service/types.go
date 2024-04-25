package service

import "time"

type Showtime struct {
	Id        int64             `json:"id"`
	BeginTime time.Time         `json:"time"`
	Duration  time.Duration     `json:"duration"`
	Title     string            `json:"title"`
	Seats     int               `json:"seats"`
	bookers   map[string]string // ip -> eventId
}
