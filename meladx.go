// ========================================
// go_meladx
// re ecriture en go du binaire meladx
// ========================================

// =============================================
// Arguments :
// -v : mode verbeux
// -r : expediteur  : adresse mail
// -s : serveur stmp  : adresse ip ou fqdn
// =============================================

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"runtime"
	"github.com/BurntSushi/toml"
	"net/smtp"
)

// Structure pour les parametres
type parametre struct {
	verbose     bool
	smtp_server string
	sender      string
}

// Structure pour le fichier de config
// Attention Majuscule obligatoire
type Config_File struct {
	Server_smtp 	string
	Auth_Login		string
	Auth_Password	string
}

var Param parametre

var DEBUG = true

// La variable qui contient les 
// les valeurs 
var Config_Auth Config_File

// Nom par defaut du fichier de config
func set_default_config_file () (filename string) {
	switch os := runtime.GOOS; os {
		case "linux":
			filename = "meladx.conf"
		case "windows":
			filename = "meladx.conf"
		default:
			filename = "meladx.conf"
	}
	return
}

func main() {

	// -----------------------------------------------
	// Analyse des arguments de la ligne de commande
	// -----------------------------------------------
	verbose := flag.Bool("v", false, "Mode verbeux")
	smtp_server := flag.String("s", "", "SMTP Server Name")
	sender := flag.String("r", "", "Sender adresse Name")

	flag.Parse()

	Param.verbose = *verbose
	Param.smtp_server = *smtp_server
	Param.sender = *sender

	if DEBUG {
	 	log.Println("========== START CMD LINE ARGS  =================" )
		log.Println(fmt.Sprintf(" Verbose Mode : %t ", Param.verbose))
		log.Println(fmt.Sprintf(" Smtp Server  : %s ", Param.smtp_server))
		log.Println(fmt.Sprintf(" Expediteur   : %s ", Param.sender))
	}

	// ------------------------------------------
	// Lecture du fichier de config par defaut
	// ------------------------------------------
	 if _, err := toml.DecodeFile(set_default_config_file(), &Config_Auth); err != nil {
				 fmt.Println("Err = %s " , err )
	 }

	 if DEBUG {
	 	log.Println("========== START CONFIG FILE =================" )
		log.Println(" Server_smtp ", Config_Auth.Server_smtp )
		log.Println(" Auth_Login ", Config_Auth.Auth_Login )
		log.Println(" Auth_Password ", Config_Auth.Auth_Password )
	 }

	// ---------------------------------
	// Analyse de l'entree standard
	// ---------------------------------
	s := bufio.NewScanner(os.Stdin)

	// TODO : detection de l'abscence de STDIN

	FROM := []string{}
	TO := []string{}
	CC := []string{}
	SUBJECT := []string{}
	PJ := []string{}
	BODY := []string{}

	// Prise en compte de la ligne de commande
	FROM = append(FROM, Param.sender)

	for s.Scan() {
		// From
		if strings.HasPrefix(s.Text(), "From:") {
			re := regexp.MustCompile(`From:(.*)$`)
			f := re.FindStringSubmatch(s.Text())
			FROM = append(FROM, strings.TrimSpace(f[1]))
			// TO
		} else if strings.HasPrefix(s.Text(), "To:") {
			re := regexp.MustCompile(`To:(.*)$`)
			f := re.FindStringSubmatch(s.Text())
			TO = append(TO, strings.TrimSpace(f[1]))
			// CC
		} else if strings.HasPrefix(s.Text(), "Cc:") {
			re := regexp.MustCompile(`Cc:(.*)$`)
			f := re.FindStringSubmatch(s.Text())
			CC = append(CC, strings.TrimSpace(f[1]))
			// SUBJECT
		} else if strings.HasPrefix(s.Text(), "Subject:") {
			re := regexp.MustCompile(`Subject:(.*)$`)
			f := re.FindStringSubmatch(s.Text())
			SUBJECT = append(SUBJECT, strings.TrimSpace(f[1]))
			// # Pieces jointes
		} else if strings.HasPrefix(s.Text(), "#") {
			PJ = append(PJ, strings.TrimSpace(s.Text()))
			// Le reste c'est in the body
		} else {
			BODY = append(BODY, s.Text())
		}
	}

	// TODO : traitement des Pieces jointes

	// Affichage du DEBUG
	if DEBUG {
		log.Println("====FROM============================")
		for _, F := range FROM {
			log.Println(F)
		}
		log.Println("====TO==============================")
		for _, F := range TO {
			log.Println(F)
		}
		log.Println("=====CC=============================")
		for _, F := range CC {
			log.Println(F)
		}
		log.Println("======SUBJECT=======================")
		for _, F := range SUBJECT {
			log.Println(F)
		}
		log.Println("=====PJ=============================")
		for _, F := range PJ {
			log.Println(F)
		}
		log.Println("=====BODY===========================")
		for _, F := range BODY {
			log.Println(F)
		}
	}

	// ---------------------
	// Envoi du mail
	// ---------------------
	if DEBUG {
		log.Println("===== ENVOI DU MAIL ===========================")
	}

	// Set up authentication information.
	auth := smtp.PlainAuth("", Config_Auth.Auth_Login, Config_Auth.Auth_Password, Config_Auth.Server_smtp )

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to      := []string{ TO[0] }
	msg_to  := fmt.Sprintf("To: %s\r\n", TO[0] )
	subject := fmt.Sprintf( "Subject: %s\r\n\n" , SUBJECT[0] )
	body    := strings.Join( BODY, "\n" )

	if DEBUG {
		log.Println(" TO : ", to )
		log.Println(" SUBJECT : ", subject )
		log.Println(" BODY : ", body )
	}

	// Construction du message
	msg := []byte( msg_to + subject + body + "\r\n" )

	err := smtp.SendMail( Config_Auth.Server_smtp, auth, Param.sender, to, msg )

	if err != nil {
		log.Fatal(err)
	}
}
