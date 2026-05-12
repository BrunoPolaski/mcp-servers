package sessions

import "time"

const (
	CookieName      = "sid"
	IdleTimeout     = 2 * time.Hour
	AbsoluteTimeout = 24 * time.Hour
)
