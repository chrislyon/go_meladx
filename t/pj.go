package main

import "os"
import "fmt"
import "bufio"
import "bytes"
import "regexp"
import "io/ioutil"
import "encoding/base64"

type PJ struct {
	type_mime string 
	description string
	filename string
	body string
}


func main(){

	f,_  := os.Open("pj.txt")
	defer f.Close()

	s := bufio.NewScanner(f)

	var re_pj = regexp.MustCompile(`\# (?P<type_mime>.*) \[(?P<description>.*)\] \"(?P<filename>.*)\"$`)

	for s.Scan() {

		pj := s.Text()
		fmt.Println("=>", pj)

		result := re_pj.FindStringSubmatch(pj)

		fmt.Println(" type_mime = ", result[1])
		fmt.Println(" Desc      = ", result[2])
		fmt.Println(" Filename  = ", result[3])

		p := PJ{ result[1], result[2], result[3], "" }

		content, _ := ioutil.ReadFile( p.filename )
		data64 := base64.StdEncoding.EncodeToString(content)

		// Decoupe en lignes
		var buf bytes.Buffer

		l_max := 500
		nb_lines := len(data64) / l_max

		for i := 0 ; i < nb_lines ; i++ {
			buf.WriteString(data64[i*l_max:(i+1)*l_max]+"\n")
		}

		buf.WriteString(data64[nb_lines*l_max:])

		p.body = buf.String()
	
		fmt.Println("p=", p, len(p.body))

	}
}
