/*
Clock Server is a concurrent TCP server that writes the time of a given timezone.
*/
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

// Connects to server and writes the timezone into it
func handleConn(c net.Conn, timeOnTimezone string) {
	defer c.Close()
	_, err := io.WriteString(c, timeOnTimezone)
	if err != nil {
		return // e.g., client disconnected
	}
}

// Reads the flag and enviroment variable values
func readParameters() (string, int) {
	timezone := os.Getenv("TZ")                      // TZ parameter
	var port = flag.Int("port", 9000, "port number") // port parameter
	flag.Parse()
	return timezone, *port
}

// TimeIn returns the time in UTC if the name is "" or "UTC"
// It returns the local time if the name is "Local".
// Otherwise, the name is taken to be a location name in
// the IANA Time Zone database, such as "Africa/Lagos".
func TimeIn(t time.Time, name string) (time.Time, error) {
	loc, err := time.LoadLocation(name)
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}

// Returns a formatted string with the given timezone and its time
// by using the TimeIn function
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

// main goroutine
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
