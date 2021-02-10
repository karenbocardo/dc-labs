package main

import (
	"fmt"
	"os"
)

func main() {

  if len(os.Args) < 2 {
        fmt.Println("Error")
    } else {

        //name := os.Args[1]

        name := os.Args[1]
        if len(os.Args) > 2 {
          for _,word := range os.Args[2:]{
            name = fmt.Sprintf("%v %v", name, word)
          }
        }

        fmt.Println("Hello " + name + ", Welcome to the Jungle")
    }

}
