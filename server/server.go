package main

import (
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "os"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, // Allow all connections
}

func connectDB() error {
    dbUser := os.Getenv("USERNAME")
    dbPassword := os.Getenv("PASSWORD")
    dbHost := os.Getenv("HOST")
    dbName := os.Getenv("DATABASE")

    dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    defer db.Close()
    fmt.Println("Database connection successful!")
    return nil
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
    // Upgrade initial GET request to a WebSocket
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer ws.Close()

    for {
        _, msg, err := ws.ReadMessage()
        if err != nil {
            fmt.Println("read error:", err)
            break
        }
        fmt.Printf("Received: %s\n", msg)

        if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
            fmt.Println("write error:", err)
            break
        }
    }
}

func main() {
    if err := connectDB(); err != nil {
        fmt.Println(err)
        return
    }
    http.HandleFunc("/ws", handleConnections)

    fmt.Println("WebSocket server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        fmt.Println("ListenAndServe:", err)
    }
}