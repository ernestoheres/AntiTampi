package main

import (
	"io"
	"crypto/sha256"
	"os"
	"fmt"
    "github.com/gorilla/websocket"
	"net/url"
)
func connectToSocket(host string, port string) (*websocket.Conn, error) {
    endpoint := "/ws"
    u := url.URL{Scheme: "ws", Host: host + ":" + port, Path: endpoint}

    conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        return nil, err
    }
    return conn, nil
}

func disconnectFromSocket(conn *websocket.Conn) error {
    err := conn.Close()
    if err != nil {
        return err
    }
    return nil
}


func main() {
	filesToHash := []string{"bla/1.txt", "bla/2.txt", "bla/3.txt"}
	correctHashes := []string{"daf26f689777e1885386c3c0c7f676b1cc0cd8208efe2d708b66afdba5a215f0","252f10c83610ebca1a059c0bae8255eba2f95be4d1d7bcfa89d7248a82d9f111","252f10c83610ebca1a059c0bae8255eba2f95be4d1d7bcfa89d7248a82d9f111" }
	combinedHash := ""
	for i,file := range filesToHash {
		hasher := sha256.New()
		f, err := os.Open(file)
		io.Copy(hasher, f)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		// print the hash
		hashStr := fmt.Sprintf("%x", hasher.Sum(nil))
		if hashStr != correctHashes[i] {
			fmt.Println("Hash does not match")
			fmt.Printf("Hash: %s\n %s \n %d \n", hashStr, correctHashes[i], i)
		}
		combinedHash += hashStr
	}	
	fmt.Printf("Combined hash: %x\n", sha256.Sum256([]byte(combinedHash)))
	conn, err := connectToSocket("localhost", "8080")
	if err != nil {
		panic(err)
	}
	defer disconnectFromSocket(conn)
	err = conn.WriteMessage(websocket.BinaryMessage, []byte(combinedHash))
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent hash to server")
}