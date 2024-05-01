package Event



type Event struct {
	ID   				string		`json:"id"`
	Name 				string		`json:"name"`
	Date 				string		`json:"date"`
	TotalTickets		int			`json:"totalTickets"`
	AvailableTickets	int			`json:"availableTickets"`
}