package main

import (
	"fmt"
	"time"
)

type Position struct {
	when      time.Time
	latitude  float64
	longitude float64
}

type Pace struct {
	mins int64
	secs int64
}

func (p *Pace) String() string {
	return fmt.Sprintf("%d:%02d", p.mins, p.secs)
}
