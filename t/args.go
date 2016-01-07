
package main

import (
	"fmt"
	"flag"
	"os"
)


func main() {
	fmt.Println(len(os.Args), os.Args)

	verbose := flag.Bool("v", false, "Mode verbeux")
	smtp_server := flag.String("s", "", "SMTP Server Name")
	sender := flag.String("r", "", "Sender adresse Name")
	config_file := flag.String("f", "", "Config File Name")

	flag.Parse()

	fmt.Printf(" Verbose Mode : %s \n", *verbose )
	fmt.Printf(" Smtp Server  : %s \n", *smtp_server )
	fmt.Printf(" Expediteur   : %s \n", *sender )
	fmt.Printf(" Config File  : %s \n", *config_file )

}
