package main

import "os"
import "fmt"
import "bufio"
import "regexp"
//import "strings"

func main(){

	f,_  := os.Open("pj.txt")
	s := bufio.NewScanner(f)

	var re_pj = regexp.MustCompile(`\# (?P<type_mime>.*) \[(?P<description>.*)\] \"(?P<filename>.*)\"$`)

	for s.Scan() {

		pj := s.Text()
		fmt.Println("=>", pj)

		result := re_pj.FindStringSubmatch(pj)

		fmt.Println(" type_mime = ", result[1])
		fmt.Println(" Desc      = ", result[2])
		fmt.Println(" Filename  = ", result[3])
	}
}
