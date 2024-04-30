package server

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"strconv"
	"strings"
)

func (ts *TicketService) CreateEvent(name string, date time.Time, totalTickets int) (*Event, error) {
	// Create a new event
	event := &Event{
		ID:               generateUUID(),
		Name:             name,
		Date:             date,
		TotalTickets:     totalTickets,
		AvailableTickets: totalTickets,
	}

	// CHECKME: IS it correct?
	go func() {
	ts.events.Store(event.ID, event)
	}()

	fmt.Println("event created", event.ID)

	return event, nil
}

func generateUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return ""
	}
	uuid[8] = (uuid[8] | 0x80) & 0xBF
	uuid[6] = (uuid[6] | 0x40) & 0x4F
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func (ts *TicketService) BookTickets(eventID string, numTickets int) ([]string, error) {
	ts.mutex.Lock()
    defer ts.mutex.Unlock()

	event, ok := ts.events.Load(eventID)
	
	if !ok {
		return nil, fmt.Errorf("event not found")
	}

	ev := event.(*Event)
	if ev.AvailableTickets < numTickets {
		return nil, fmt.Errorf("not enough tickets available")
	}

	var ticketIDs []string
	for i := 0; i < numTickets; i++ {
		ticket := &Ticket{
			ID:      generateUUID(),
			EventID: eventID,
		}
		ts.tickets.Store(ticket.ID, ticket)
		ticketIDs = append(ticketIDs, ticket.ID)
	}

	ev.AvailableTickets -= numTickets
	ts.events.Store(eventID, ev)

	return ticketIDs, nil
}

func (ts *TicketService) ListEvents() []*Event {
	var events []*Event
	ts.events.Range(func(key, value interface{}) bool {
		event := value.(*Event)
		events = append(events, event)
		return true
	})
	return events
}

func (ts *TicketService) getListEventsHandler(w http.ResponseWriter, r *http.Request) {
	events := ts.ListEvents()

	fmt.Fprintf(w, "{\"events\": [")
	for i, event := range events {
		fmt.Fprintf(w, "{\"id\": \"%s\", \"name\": \"%s\", \"date\": \"%s\", \"total_tickets\": %d, \"available_tickets\": %d}", event.ID, event.Name, event.Date.Format("2006-01-02"), event.TotalTickets, event.AvailableTickets)
		if i < len(events)-1 {
			fmt.Fprintf(w, ", ")
		}
	}
	fmt.Fprintf(w, "]}")
}

func getHomePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the ticket service")
}

func (ts *TicketService) reserveTicketsHandler(w http.ResponseWriter, r *http.Request) {

	parsedURL, err := url.Parse(r.URL.String())
	pathParts := strings.Split(parsedURL.Path, "/")
	eventID := pathParts[2]
	numTickets, err := strconv.Atoi(r.URL.Query().Get("num_tickets"))

	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		return
	}

	ticketIDs, err := ts.BookTickets(eventID, numTickets)
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		return
	}

	fmt.Fprintf(w, "{\"tickets\": {\"event_id\": \"%s\", \"ticket_ids\": [", eventID)
	for i, ticketID := range ticketIDs {
		fmt.Fprintf(w, "\"%s\"", ticketID)
		if i < len(ticketIDs)-1 {
			fmt.Fprintf(w, ", ")
		}
	}
	fmt.Fprintf(w, "]}}")
}

func (ts *TicketService) createEventHandler(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	date, err := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		return
	}
	totalTickets, err := strconv.Atoi(r.URL.Query().Get("total_tickets"))
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		return
	}
	event, err := ts.CreateEvent(name, date, totalTickets)
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		return
	}
	fmt.Fprintf(w, "{\"event\": {\"id\": \"%s\", \"name\": \"%s\", \"date\": \"%s\", \"total_tickets\": %d, \"available_tickets\": %d}}", event.ID, event.Name, event.Date.Format("2006-01-02"), event.TotalTickets, event.AvailableTickets)
}