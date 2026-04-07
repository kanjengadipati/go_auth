package main

import (
	"log"

	"go-auth-app/config"
	"go-auth-app/seeds"
)

func main() {
	// Load env (WAJIB)
	config.LoadEnv()

	// Init DB (WAJIB)
	config.ConnectDB()

	// Run seeder
	log.Println("Start seeding...")

	seeds.SeedRoles(config.DB)
	log.Println("SeedRoles done")

	seeds.SeedPermissions(config.DB)
	log.Println("SeedPermissions done")

	seeds.SeedAdmin(config.DB)
	log.Println("SeedAdmin done")

	log.Println("Seeding done 🚀")
}
