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

// Librairie necessaire pour la prise en compte
// des fichiers de configs
// go get github.com/BurntSushi/toml

package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
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
	Port			int
	Auth_Login		string
	Auth_Password	string
}

// Structure pour les pieces jointes
type PJ struct {
	type_mime string 
	description string
	filepath string
	filename string
	body string
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
			filename = "/etc/meladx.conf"
			filename = "meladx.conf"
		case "windows":
			filename = "meladx.conf"
		default:
			filename = "meladx.conf"
	}
	return
}

func f_path( p string) (filename string) {
	switch os := runtime.GOOS; os {
		case "windows":
			t := strings.Split(p, "\\")
			filename = t[len(t)-1]
		default:
			filename = path.Base(p)
	}
	return
}

func encode_pj( pj string )( PJ ) {
	
	var re_pj = regexp.MustCompile(`\# (?P<type_mime>.*) \[(?P<description>.*)\] \"(?P<filename>.*)\"$`)

	//fmt.Println("=>", pj)

	result := re_pj.FindStringSubmatch(pj)

	//fmt.Println(" type_mime = ", result[1])
	//fmt.Println(" Desc      = ", result[2])
	//fmt.Println(" Filename  = ", result[3])

	p := PJ{ result[1], result[2], result[3], f_path(result[3]), "" }

	content, err := ioutil.ReadFile( p.filepath )

	if err == nil {

		data64 := base64.StdEncoding.EncodeToString(content)

		// Decoupe en lignes
		var buf bytes.Buffer

		l_max := 80
		nb_lines := len(data64) / l_max

		for i := 0 ; i < nb_lines ; i++ {
			buf.WriteString(data64[i*l_max:(i+1)*l_max]+"\n")
		}

		buf.WriteString(data64[nb_lines*l_max:])

		p.body = buf.String()

		//fmt.Println("p=", p, len(p.body))
	}

	return p

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
		log.Println(" Config File   :", set_default_config_file() )
		log.Println(" Server_smtp   :", Config_Auth.Server_smtp )
		log.Println(" Port          :", Config_Auth.Port )
		log.Println(" Auth_Login    :", Config_Auth.Auth_Login )
		log.Println(" Auth_Password :", Config_Auth.Auth_Password )
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
			t := strings.Split(f[1], ";")
			for _, l := range t {
				TO = append(TO, strings.TrimSpace(l))
			}
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


	m := gomail.NewMessage()


	//From (expediteur)
	m.SetHeader( "From", FROM[0] )

	// La liste des destinataires
	for _, t := range TO {
		m.SetHeader( "To", t )
	}

	// La liste des destinataires CC
	for _, t := range CC {
		m.SetHeader( "Cc", t )
	}

	// Le sujet
	m.SetHeader( "Subject", fmt.Sprintf("%s", SUBJECT[0] ))

	// Le Corps
	m.SetBody( "text/html", strings.Join( BODY, "\n" ))

	// Les pieces jointes
	log.Println(" Pieces jointes : ")
	for _, pj := range PJ {
		log.Println("PJ : ", pj)
		//p := encode_pj(pj)
		//m.Attach(p.filepath)
	}

	fmt.Println("Message : ", m)

	// Envoi du message

	d := gomail.NewPlainDialer( Config_Auth.Server_smtp, Config_Auth.Port, Config_Auth.Auth_Login, Config_Auth.Auth_Password )

	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Email envoyé", err)
	}
}
