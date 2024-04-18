package gomessage_test

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/jaxoncarelos/GoMessage"
)

func TestAddCommand(t *testing.T) {
	server := gomessage.NewGoMessage()
	store := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	log.SetOutput(w)
	log.SetFlags(0)
	server.OnMessage = func(str string, conn net.Conn) {
	}
	server.OnConnect = func(conn net.Conn) {
	}
	server.AddCommand("p", func(str string) error {
		log.Println("Hello from command: ", str)
		return nil
	})
	go server.Start(":8080")
	// create a tcp connection to 8080 and send a dummy message
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	conn.Write([]byte("p:12:Hello World!"))
	// wait 3 seconds and stop
	time.Sleep(3 * time.Second)
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = store

	fmt.Println("We got: ", string(out))
	fmt.Println("Checking if it contains: \"Hello from command: 12:Hello World!\"")
	if string(out) != "Hello from command:  12:Hello World!\n" {
		t.Fail()
	}
	conn.Close()
}
