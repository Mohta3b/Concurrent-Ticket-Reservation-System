package TicketService

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	// "log"
	"os"
	"sync"
	"ticket_reservation/src/Event"
	"ticket_reservation/src/Ticket"
	"time"
)

type TicketService struct {
	events  sync.Map
	tickets sync.Map
	mutex   sync.RWMutex
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

func (ts *TicketService) LoadEvents() error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	eventsFilePath := "./data/events.json"
	eventsFile, err := os.Open(eventsFilePath)
	if err != nil {
		return fmt.Errorf("error opening events file: %v", err)
	}
	defer eventsFile.Close()

	var events []*Event.Event
	err = json.NewDecoder(eventsFile).Decode(&events)
	if err != nil {
		return fmt.Errorf("error decoding events file: %v", err)
	}

	for _, event := range events {
		ts.events.Store(event.ID, event)
	}

	return nil
}

func (ts *TicketService) LoadTickets() error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()

	ticketsFilePath := "./data/tickets.json"
	ticketsFile, err := os.Open(ticketsFilePath)
	if err != nil {
		return fmt.Errorf("error opening tickets file: %v", err)
	}
	defer ticketsFile.Close()

	var tickets []*Ticket.Ticket
	err = json.NewDecoder(ticketsFile).Decode(&tickets)
	if err != nil {
		return fmt.Errorf("error decoding tickets file: %v", err)
	}

	for _, ticket := range tickets {
		ts.tickets.Store(ticket.ID, ticket)
	}

	return nil
}

func (ts *TicketService) CreateEvent(name string, date time.Time, totalTickets int) (*Event.Event, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	event := &Event.Event{
		ID:               generateUUID(),
		Name:             name,
		Date:             date.Format("2006-01-02"),
		TotalTickets:     totalTickets,
		AvailableTickets: totalTickets,
	}

	// CHECKME: IS it correct?
	ts.events.Store(event.ID, event)

	return event, nil
}

func (ts *TicketService) ListEvents() []*Event.Event {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	var events []*Event.Event
	ts.events.Range(func(key, value interface{}) bool {
		event, ok := value.(*Event.Event)
		if ok {
			events = append(events, event)
		}
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
		return nil, fmt.Errorf("not enough tickets available for event %s -> %s (availableTickets: %d , requestedTickets: %d)", eventID, ev.Name, ev.AvailableTickets, numTickets)
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

func (ts *TicketService) GetEvent(eventID string) *Event.Event {
	event, ok := ts.events.Load(eventID)
	if !ok {
		return nil
	}
	return event.(*Event.Event)
}
