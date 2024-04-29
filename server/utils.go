package server

import (
	"time"
	"crypto/rand"
	"fmt"
)


func (ts *TicketService) CreateEvent(name string, date time.Time, totalTickets int) (*Event, error) {
	// Create a new event
	event := &Event{
		ID:   				generateUUID(),
		Name: 				name,
		Date: 				date,
		TotalTickets:		totalTickets,
		AvailableTickets:	totalTickets,
	}
	
	// Save the event
	ts.events.Store(event.ID, event)
	
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
	// implement concurrency control here 
	// ...

	event, ok := ts.events.Load(eventID)
	if !ok {
		return nil, fmt.Errorf("event not found")
	}

	// check if there are enough tickets available
	ev := event.(*Event)
	if ev.AvailableTickets < numTickets {
		return nil, fmt.Errorf("not enough tickets available")
	}

	// create tickets
	var ticketIDs []string
	for i := 0; i < numTickets; i++ {
		ticket := &Ticket{
			ID:     generateUUID(),
			EventID: eventID,
		}
		ts.tickets.Store(ticket.ID, ticket)
		ticketIDs = append(ticketIDs, ticket.ID)
	}

	// update available tickets
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