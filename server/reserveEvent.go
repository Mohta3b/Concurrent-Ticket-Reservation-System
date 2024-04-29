package server

import (
	"fmt"
)

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