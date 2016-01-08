//
// go_meladx
// re ecriture en go du binaire meladx
//

// Arguments :
// -v : mode verbeux
// -r : expediteur  : adresse mail
// -s : serveur stmp  : adresse ip ou fqdn
// -f : config file (a voir)

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

// Structure pour les parametres
type parametre struct {
	verbose     bool
	smtp_server string
	sender      string
	config_file string
}

var Param parametre

func main() {
	fmt.Println(len(os.Args), os.Args)

	// -----------------------------------------------
	// Analyse des arguments de la ligne de commande
	// -----------------------------------------------
	verbose := flag.Bool("v", false, "Mode verbeux")
	smtp_server := flag.String("s", "", "SMTP Server Name")
	sender := flag.String("r", "", "Sender adresse Name")
	config_file := flag.String("f", "", "Config File Name")

	flag.Parse()

	Param.verbose = *verbose
	Param.smtp_server = *smtp_server
	Param.sender = *sender
	Param.config_file = *config_file

	fmt.Printf(" Verbose Mode : %t \n", Param.verbose)
	fmt.Printf(" Smtp Server  : %s \n", Param.smtp_server)
	fmt.Printf(" Expediteur   : %s \n", Param.sender)
	fmt.Printf(" Config File  : %s \n", Param.config_file)

	// Analyse de l'entree standard
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		log.Println("line : ", s.Text())
	}

}
