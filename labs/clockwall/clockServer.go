/*
This program must accept the -port parameter
TZ=US/Eastern    go run serverClock.go -port 8010
*/

// Clock Server is a concurrent TCP server that periodically writes the time.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func handleConn(c net.Conn, timeOnTimezone string) {
	defer c.Close()
	_, err := io.WriteString(c, timeOnTimezone)
	if err != nil {
		return // e.g., client disconnected
	}
}

func readParameters() (string, int) {
	// read parameters
	timezone := os.Getenv("TZ")                      // TZ parameter
	var port = flag.Int("port", 9000, "port number") // port parameter
	flag.Parse()
	return timezone, *port
}

// TimeIn returns the time in UTC if the name is "" or "UTC"
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

func getTime(timezone string) string {

	t, err := TimeIn(time.Now(), timezone)
	if err == nil {
		return fmt.Sprintf("%v\t: %v\n", t.Location(), t.Format("15:04:05"))
		// fmt.Println(t.Location(), "\t:", t.Format("15:04:05"))
	} else {
		return fmt.Sprintf("%v\t: <time unknown>\n", t.Location())
		// fmt.Println(timezone, "<time unknown>")
	}
}

func main() {
	// read parameters
	timezone, port := readParameters()
	// get time for timezone
	timeOnTimezone := getTime(timezone)

	// fmt.Printf("Timezone: %v\n", timezone)
	// fmt.Printf("Port: %v\n", port)
	// fmt.Println(timeOnTimezone)

	server := fmt.Sprint("localhost:", port)
	// fmt.Println(site)

	// conection
	listener, err := net.Listen("tcp", server) // pass port
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn, timeOnTimezone) // handle connections concurrently
	}
}
