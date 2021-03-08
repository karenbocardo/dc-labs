/*
This program acts as a client of several clock servers at once
*/
package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func dialServer(ch chan string) {
	for v := range ch {
		// fmt.Println("read value", v, "from ch")
		server := v

		conn, err := net.Dial("tcp", server)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		mustCopy(os.Stdout, conn)
	}
}

func write(args []string, ch chan string) {
	for _, arg := range args {
		// fmt.Println(arg)
		server := strings.Split(arg, "=")[1]
		ch <- server
		// fmt.Println("successfully wrote", server, "to ch")
		// fmt.Println(server)
	}
	close(ch)
}

func main() {
	// parameters
	args := os.Args[1:]
	// fmt.Println(args)

	ch := make(chan string) // unbuffered channel
	go write(args, ch)
	dialServer(ch)
}
