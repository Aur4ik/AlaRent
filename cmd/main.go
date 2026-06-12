package main

import(
	"log"

	
	"github.com/Aur4ik/AlaRent/pkg/config"
	"github.com/Aur4ik/AlaRent/pkg/database"
)

func main(){
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db); err != nil{
		log.Fatal(err)
	}

	log.Println("DB is working")

	_ = db

}