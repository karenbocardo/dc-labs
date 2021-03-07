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

func main() {
	// parameters
	args := os.Args[1:]
	// fmt.Println(args)

	for _, arg := range args {
		// fmt.Println(arg)
		server := strings.Split(arg, "=")[1]
		//fmt.Println(server)

		conn, err := net.Dial("tcp", server)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		mustCopy(os.Stdout, conn)
	}

}

/*
unbuffered channel ?
recibir las fechas de los tres servidores al mismo tiempo, vaciarlo al imprimirlo y en ese instante se llena con otra de otro servidor


*/
