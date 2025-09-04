package model

import "net/url"

// An Organizer in the iCalendar format, used in VEVENT, VTODO, and VJOURNAL
// for more information see https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.3
type Organizer struct {
	// denoted by CN= in the spec
	CommonName string
	// Note: Any Valid URI
	// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.3
	CalAddress *url.URL
	// denoted by DIR= in the spec
	Directory string
}
