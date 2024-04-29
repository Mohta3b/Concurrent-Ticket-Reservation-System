package server

func (ts *TicketService) ListEvents() []*Event {
	var events []*Event
	ts.events.Range(func(key, value interface{}) bool {
		event := value.(*Event)
		events = append(events, event)
		return true
	})
	return events
}