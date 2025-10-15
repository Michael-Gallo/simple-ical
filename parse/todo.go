package parse

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/michael-gallo/simple-ical/model"
)

const todoLocation = "Todo"

// parseTodoProperty parses a single property line and adds it to the provided todo.
func parseTodoProperty(propertyName string, value string, params map[string]string, todo *model.Todo) error {
	switch model.TodoToken(propertyName) {
	case model.TodoTokenDTStamp:
		return setOnceTimeProperty(&todo.DTStamp, value, propertyName, todoLocation)
	case model.TodoTokenUID:
		return setOnceProperty(&todo.UID, value, propertyName, todoLocation)
	case model.TodoTokenClass:
		return setOnceProperty(&todo.Class, model.TodoClass(value), propertyName, todoLocation)
	case model.TodoTokenCompleted:
		return setOnceTimeProperty(&todo.Completed, value, propertyName, todoLocation)
	case model.TodoTokenCreated:
		return setOnceTimeProperty(&todo.Created, value, propertyName, todoLocation)
	case model.TodoTokenDescription:
		todo.Description = append(todo.Description, value)
		return nil
	case model.TodoTokenDTStart:
		return setOnceTimeProperty(&todo.DTStart, value, propertyName, todoLocation)

	// Due and Duration are mutually exclusive
	case model.TodoTokenDue:
		if todo.Duration != 0 {
			return errInvalidDurationPropertyDue
		}
		return setOnceTimeProperty(&todo.Due, value, propertyName, todoLocation)
	case model.TodoTokenDuration:
		if todo.Due != (time.Time{}) {
			return errInvalidDurationPropertyDue
		}
		return setOnceDurationProperty(&todo.Duration, value, propertyName, todoLocation)

	case model.TodoTokenGeo:
		if todo.Geo != nil {
			return fmt.Errorf("%w: %s", errDuplicateProperty, propertyName)
		}
		// Geo must be two floats separated by a semicolon
		latitudeString, longitudeString, found := strings.Cut(value, ";")
		if !found {
			return errInvalidGeoProperty
		}
		latitude, err := strconv.ParseFloat(latitudeString, 64)
		if err != nil {
			return errInvalidGeoPropertyLatitude
		}
		longitude, err := strconv.ParseFloat(longitudeString, 64)
		if err != nil {
			return errInvalidGeoPropertyLongitude
		}
		todo.Geo = append(todo.Geo, latitude, longitude)
	case model.TodoTokenLastModified:
		return setOnceTimeProperty(&todo.LastModified, value, propertyName, todoLocation)
	case model.TodoTokenLocation:
		return setOnceProperty(&todo.Location, value, propertyName, todoLocation)
	case model.TodoTokenOrganizer:
		organizer, err := parseOrganizer(value, params)
		if err != nil {
			return err
		}
		todo.Organizer = organizer
	case model.TodoTokenPercentComplete:
		return setOnceIntProperty(&todo.PercentComplete, value, propertyName, todoLocation)
	case model.TodoTokenPriority:
		return setOnceIntProperty(&todo.Priority, value, propertyName, todoLocation)
	case model.TodoTokenRecurrenceID:
		return setOnceTimeProperty(&todo.RecurrenceID, value, propertyName, todoLocation)
	case model.TodoTokenSequence:
		return setOnceIntProperty(&todo.Sequence, value, propertyName, todoLocation)
	case model.TodoTokenStatus:
		todo.Status = model.TodoStatus(value)
	case model.TodoTokenSummary:
		return setOnceProperty(&todo.Summary, value, propertyName, todoLocation)
	case model.TodoTokenTransp:
		return setOnceProperty(&todo.Transp, model.TodoTransp(value), propertyName, todoLocation)
	case model.TodoTokenURL:
		return setOnceProperty(&todo.URL, value, propertyName, todoLocation)

	// Repeatable properties
	case model.TodoTokenAttach:
		todo.Attach = append(todo.Attach, value)
		return nil
	case model.TodoTokenAttendee:
		parsedURL, err := url.Parse(value)
		if err != nil {
			return err
		}
		todo.Attendees = append(todo.Attendees, *parsedURL)
	case model.TodoTokenCategories:
		todo.Categories = append(todo.Categories, strings.Split(value, ",")...)
	case model.TodoTokenComment:
		todo.Comment = append(todo.Comment, value)
	case model.TodoTokenContact:
		todo.Contacts = append(todo.Contacts, value)
	case model.TodoTokenExceptionDates:
		return appendTimeProperty(&todo.ExceptionDates, value, propertyName, todoLocation)
	case model.TodoTokenRequestStatus:
		todo.RequestStatus = append(todo.RequestStatus, value)
	case model.TodoTokenRelated:
		todo.Related = append(todo.Related, value)
	case model.TodoTokenResources:
		todo.Resources = append(todo.Resources, strings.Split(value, ",")...)
	case model.TodoTokenRdate:
		return appendTimeProperty(&todo.Rdate, value, propertyName, todoLocation)
	default:
		return fmt.Errorf("%w: %s", errInvalidTodoProperty, propertyName)
	}
	return nil
}

// validateTodo ensures that all required values are present for a todo.
func validateTodo(ctx *parseContext) error {
	if ctx.currentTodo.UID == "" {
		return errMissingTodoUIDProperty
	}
	if ctx.currentTodo.DTStart == (time.Time{}) {
		return errMissingTodoDTStartProperty
	}
	return nil
}
