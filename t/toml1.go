package main

import "github.com/BurntSushi/toml"
import "fmt"

type Config struct {
	Server_smtp string
	Auth_Login	string
	Auth_Password string
}

func main() {

	var config Config
	var f string
	f = "t1.toml"
	if _, err := toml.DecodeFile(f, &config); err != nil {
		fmt.Println("Err = %s " , err )
	}


	fmt.Println(" Server SMTP ", config.Server_smtp )
	fmt.Println(" User        ", config.Auth_Login )
	fmt.Println(" Password    ", config.Auth_Password )
}
