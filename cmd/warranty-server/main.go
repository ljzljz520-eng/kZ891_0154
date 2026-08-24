package main

import (
	"context"
	"log"
	"net/http"

	"warrantyservice/internal/adapters"
	"warrantyservice/internal/config"
	"warrantyservice/internal/httpapi"
	"warrantyservice/internal/model"
	"warrantyservice/internal/persistence"
	"warrantyservice/internal/support"
	"warrantyservice/internal/warranty"
)

func buildHandler(cfg config.Config) (*httpapi.Handler, *persistence.Store, error) {
	if err := cfg.Validate(); err != nil {
		return nil, nil, err
	}
	store, err := persistence.Open(cfg.DatabasePath)
	if err != nil {
		return nil, nil, err
	}
	if err := seed(context.Background(), store, cfg); err != nil {
		_ = store.Close()
		return nil, nil, err
	}
	reader := adapters.NewWarrantyAdapter(store)
	points := adapters.NewPointDirectory(store)
	audit := adapters.NewAuditWriter(store)
	advisor := support.NewAdvisor(store, cfg.SupportHotline, "官方小程序、电话、授权服务网点", cfg.FixedToday)
	service := warranty.NewService(reader, points, audit, advisor, cfg.FixedToday)
	return httpapi.NewHandler(service, cfg), store, nil
}

func seed(ctx context.Context, store *persistence.Store, cfg config.Config) error {
	points := config.FixturePoints()
	for _, record := range config.FixtureWarranties() {
		selected := make([]model.ServicePoint, 0, len(record.ServicePoints))
		for _, name := range record.ServicePoints {
			for _, point := range points {
				if point.Name == name {
					selected = append(selected, point)
				}
			}
		}
		if err := store.SaveWarranty(ctx, record, selected); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	cfg := config.Load()
	handler, store, err := buildHandler(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	log.Printf("%s listening on %s", cfg.ServiceName, cfg.Address)
	if err := http.ListenAndServe(cfg.Address, handler); err != nil {
		log.Fatal(err)
	}
}
