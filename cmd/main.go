package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mongo-db-connector/domain"
	"mongo-db-connector/internal/database"
	"mongo-db-connector/internal/infrastructure"
	"mongo-db-connector/internal/repository"
)

func main() {
	ctx := context.Background()
	conf := &infrastructure.DBConfig{
		Host:     "localhost",
		Port:     27017,
		Database: "connector",
		Username: "admin",
		Password: "admin",
	}

	dbClient, err := database.BuildClient(ctx, conf)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = dbClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	migrationConfig := &infrastructure.MigrationConfig{
		DBConfig: infrastructure.DBConfig{
			Host:     "localhost",
			Port:     27017,
			Database: "connector",
			Username: "admin",
			Password: "admin",
		},
		CollectionName: "migrations",
	}
	migrator := database.NewMigrationRunner(dbClient, migrationConfig)
	err = migrator.RunMigrations()
	if err != nil {
		panic(err)
	}

	offerRepository := repository.NewOfferRepository(dbClient, conf.Database)
	offer := &repository.DBOffer{
		ID: primitive.NewObjectID(),
		Offer: &domain.Offer{
			OfferID:     "OFFER_ID123",
			Description: "Offer 123",
		},
	}

	err = offerRepository.CreateEntity(ctx, offer)
	if err != nil {
		panic(err)
	}

	loadedOffer, err := offerRepository.GetEntity(ctx, "OFFER_ID123")
	if err != nil {
		panic(err)
	}
	fmt.Printf("OFFER LOADED: %v\n", loadedOffer)

	offer.Description = "Test offer update"
	updatedOffer, err := offerRepository.UpdateEntity(ctx, offer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("UPDATED OFFER: %v\n", updatedOffer)

	err = offerRepository.DeleteEntity(ctx, "OFFER_ID123")
	if err != nil {
		panic(err)
	}

	fmt.Println("ALL DONE")

}
