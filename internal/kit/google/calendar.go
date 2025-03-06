package google

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const primaryCalendar = "primary"

var (
	ErrEmptyCredentials = errors.New("client id or client secret cannot be empty")
)

type GoogleCalendar struct {
	service *calendar.Service
}

func NewCalendar(ctx context.Context, clientID, clientSecret string, t *oauth2.Token) (GoogleCalendar, error) {
	if clientID == "" || clientSecret == "" {
		return GoogleCalendar{}, ErrEmptyCredentials
	}

	c := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			calendar.CalendarEventsScope,
		},
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	client := c.Client(ctx, t)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return GoogleCalendar{}, fmt.Errorf("creating google calendar service: %w", err)
	}

	return GoogleCalendar{service: srv}, nil
}

type CalendarEvent struct {
	ID        string
	Summary   string
	Start     time.Time
	End       time.Time
	Attendees []string
}

type Meet struct {
	ID  string
	URL string
}

// DailyEventsUTC retrieves the events from the primary calendar within the specified day
// (fromTime to fromTime + 1 day) and converts them to UTC.
func (gc GoogleCalendar) DailyEventsUTC(dateTime time.Time) ([]CalendarEvent, error) {
	toTime := dateTime.AddDate(0, 0, 1) // + 1 day

	events, err := gc.service.Events.List(primaryCalendar).
		SingleEvents(true).
		OrderBy("startTime").
		TimeMin(dateTime.Format(time.RFC3339)).
		TimeMax(toTime.Format(time.RFC3339)).
		Do()
	if err != nil {
		return nil, fmt.Errorf("retrieving events from calendar: %w", err)
	}

	if len(events.Items) == 0 {
		return nil, nil
	}

	dailyEvents := make([]CalendarEvent, 0, len(events.Items))
	for _, event := range events.Items {
		// Some events doesn't have start/end time, like work location, we should ignore them.
		if event.Start.DateTime == "" || event.End.DateTime == "" {
			continue
		}

		fromTime, err := time.Parse(time.RFC3339, event.Start.DateTime)
		if err != nil {
			return nil, fmt.Errorf("parsing fromTime: %w", err)
		}

		toTime, err := time.Parse(time.RFC3339, event.End.DateTime)
		if err != nil {
			return nil, fmt.Errorf("parsing toTime: %w", err)
		}

		dailyEvents = append(dailyEvents, CalendarEvent{
			Summary: event.Summary,
			Start:   fromTime,
			End:     toTime,
		})
	}

	return dailyEvents, nil
}

func (gc GoogleCalendar) CreateEvent(e CalendarEvent) (Meet, error) {
	attendees := make([]*calendar.EventAttendee, 0, len(e.Attendees))
	for _, attendee := range e.Attendees {
		attendees = append(attendees, &calendar.EventAttendee{Email: attendee})
	}
	event := calendar.Event{
		Summary: e.Summary,
		Start: &calendar.EventDateTime{
			DateTime: e.Start.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: e.End.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		Attendees: attendees,
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: e.ID,
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: "hangoutsMeet",
				},
			},
		},
	}

	createdEvent, err := gc.service.Events.Insert(primaryCalendar, &event).ConferenceDataVersion(1).Do()
	if err != nil {
		return Meet{}, fmt.Errorf("inserting calendar event: %w", err)
	}

	return Meet{ID: createdEvent.Id, URL: createdEvent.HangoutLink}, nil
}

func (gc GoogleCalendar) DeleteEvent(eventID string) error {
	err := gc.service.Events.Delete(primaryCalendar, eventID).Do()
	if err != nil {
		return fmt.Errorf("cannot delete event %s from calendar: %w", eventID, err)
	}

	return nil
}
