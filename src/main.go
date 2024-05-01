package main

import (
	server "ticket_reservation/src/server"
)

const PORT = ":5050"

func main() {
	server.Run(PORT)
}