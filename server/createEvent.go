package server

import (
	"time"
	"fmt"
	"crypto/rand"
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