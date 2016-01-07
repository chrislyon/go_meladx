package main

import "os"
import "log"
import "bufio"

func main(){
	s:= bufio.NewScanner(os.Stdin)

	for s.Scan(){
		log.Println("line : ", s.Text())
	}
}
