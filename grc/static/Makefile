run:
	@go run main.go

air: 
	@air

assets:
	@grc generate-assets

db/create: 
	@docker compose up db --build

db/start:
	@docker compose up db

# TODO: Add ENV port here
templ:
	@templ generate -watch --proxy="http://localhost:3000"

dev:
	make -j2 air templ