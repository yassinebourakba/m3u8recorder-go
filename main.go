package main

import (
	"context"
	"jazzine/m3u8recorder/database"
	"jazzine/m3u8recorder/datasource"
	"jazzine/m3u8recorder/libs/logging"
	"jazzine/m3u8recorder/services"
)

func main() {
	logging.Init()

	db, err := database.NewConnection()
	if err != nil {
		logging.Log().
			WithError(err).
			Error("error connecting to db")
		panic(err)
	}
	defer db.Close()

	err = database.Migrate(*db)
	if err != nil {
		logging.Log().
			WithError(err).
			Error("error migrating db")
		panic(err)
	}

	ctx := context.Background()

	dao := datasource.NewDao(db)

	services.Record(ctx, dao)
}
