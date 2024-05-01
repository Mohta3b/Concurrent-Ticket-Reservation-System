package TicketService

import (
	"ticket_reservation/src/Event"
	"ticket_reservation/src/Ticket"
	"sync"
	"fmt"
	"crypto/rand"
	"time"
)

type TicketService struct {
	events sync.Map
	tickets sync.Map
	mutex sync.Mutex
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

func (ts *TicketService) CreateEvent(name string, date time.Time, totalTickets int) (*Event.Event, error) {
	// Create a new event
	event := &Event.Event{
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

func (ts *TicketService) ListEvents() []*Event.Event {
	var events []*Event.Event
	ts.events.Range(func(key, value interface{}) bool {
		event := value.(*Event.Event)
		events = append(events, event)
		return true
	})
	return events
}

func (ts *TicketService) BookTickets(eventID string, numTickets int) ([]string, error) {
	ts.mutex.Lock()
    defer ts.mutex.Unlock()

	event, ok := ts.events.Load(eventID)
	
	if !ok {
		return nil, fmt.Errorf("event ID %s not found", eventID)
	}

	ev := event.(*Event.Event)
	if ev.AvailableTickets < numTickets {
		return nil, fmt.Errorf("not enough tickets available for event %s", eventID)
	}

	var ticketIDs []string
	for i := 0; i < numTickets; i++ {
		ticket := &Ticket.Ticket{
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