package service

import "time"

type Showtime struct {
	Id        int               `json:"id"`
	BeginTime time.Time         `json:"beginTime"`
	EndTime   time.Time         `json:"endTime"`
	Duration  int               `json:"duration"` // in minutes
	Title     string            `json:"title"`
	Seats     int               `json:"seats"`
	bookers   map[string]string // ip -> eventId
}
