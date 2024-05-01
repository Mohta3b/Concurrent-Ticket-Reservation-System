package main

const PORT = ":5050"

func main() {
    serverURL := "http://localhost" + PORT
    eventURL := "http://localhost" + PORT + "/events"
    reserveURL := "http://localhost" + PORT + "/events/"
    
    client := NewClient(serverURL, eventURL, reserveURL)
    
    ConnectToServer(client)
}
