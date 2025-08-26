package model

import "time"

type Event struct {
	Summary     string
	Description string
	Start       time.Time
	End         time.Time
	Location    string
	Organizer   *Organizer
}

type Organizer struct {
	CommonName string
	Mailto     string
}
