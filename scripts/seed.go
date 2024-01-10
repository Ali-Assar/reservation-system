package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Ali-Assar/car-rental-system/api"
	"github.com/Ali-Assar/car-rental-system/db"
	"github.com/Ali-Assar/car-rental-system/db/fixtures"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	agencyStore := db.NewMongoAgencyStore(client)
	store := &db.Store{
		User:        db.NewMongoUserStore(client),
		Reservation: db.NewMongoReservationStore(client),
		Car:         db.NewMongoCarStore(client, agencyStore),
		Agency:      db.NewMongoAgencyStore(client),
	}
	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Println("user token ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin token ->", api.CreateTokenFromUser(admin))
	agency := fixtures.AddAgency(store, "voom voom", "rome", 2, nil)
	car := fixtures.AddCar(store, "sport", "petrol", "BMW", 2017, 100, agency.ID)
	reservation := fixtures.AddReservation(store, user.ID, car.ID, time.Now(), time.Now().AddDate(0, 0, 2))
	fmt.Println("reservation- >", reservation.ID)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("fake agency %d", i)
		location := fmt.Sprintf("fake location %d", i)

		fixtures.AddAgency(store, name, location, rand.Intn(5)+1, nil)
	}
}
