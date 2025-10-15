package parse

import "github.com/michael-gallo/simple-ical/model"

// parseCalendarProperty parses a single property line and sets its value in the provided vcalendar.
func parseCalendarProperty(propertyName string, value string, _ map[string]string, calendar *model.Calendar) error {
	switch propertyName {
	case "VERSION":
		return setOnceProperty(&calendar.Version, value, propertyName, "VCALENDAR")
	case "PRODID":
		return setOnceProperty(&calendar.ProdID, value, propertyName, "VCALENDAR")
	case "CALSCALE":
		return setOnceProperty(&calendar.CalScale, value, propertyName, "VCALENDAR")
	case "METHOD":
		return setOnceProperty(&calendar.Method, value, propertyName, "VCALENDAR")
	}
	return nil
}

func validateCalendar(calendar *model.Calendar) error {
	if calendar.Version == "" {
		return ErrMissingCalendarVersionProperty
	}
	if calendar.ProdID == "" {
		return ErrMissingCalendarProdIDProperty
	}
	return nil
}
