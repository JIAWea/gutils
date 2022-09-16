package gtime

import "time"

const (
	Nanosecond  = time.Nanosecond
	Microsecond = time.Microsecond
	Millisecond = time.Millisecond
	Second      = time.Second
	Minute      = time.Minute
	HalfHour    = time.Minute * 30
	Hour        = time.Hour
	HalfDay     = time.Hour * 12
	Day         = time.Hour * 24
	Week        = Day * 7
	Month       = Day * 30
	Quarter     = Day * 91
	Year        = Day * 365
)
