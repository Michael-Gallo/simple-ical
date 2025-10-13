// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package model

import (
	"net/url"
)

// Organizer represents an ORGANIZER component in the iCalendar format, used in VEVENT, VTODO, and VJOURNAL
// for more information see https://datatracker.ietf.org/doc/html/rfc5545#section-3.8.4.3
type Organizer struct {
	// denoted by CN
	//See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.2.2
	CommonName string
	// Note: Any Valid URI
	// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.3
	CalAddress *url.URL

	// denoted by DIR
	// A directory entry reference
	// See: https://datatracker.ietf.org/doc/html/rfc5545#section-3.2.6
	Directory *url.URL

	// denoted by SENT-BY
	// See https://datatracker.ietf.org/doc/html/rfc5545#section-3.2.18
	SentBy *url.URL

	// denoted by LANGUAGE
	// https://datatracker.ietf.org/doc/html/rfc5545#section-3.2.10
	// no validation is done on the string at this time, but it is intended to be a valid tag under RFC5646
	// See: https://datatracker.ietf.org/doc/html/rfc5646
	Language string

	OtherParams map[string]string
}
