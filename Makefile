# Local dev environment setup
setup:
	./scripts/genDotenv.sh

# Run the application for local development (do not use in deployments)
run:
	go run application.go

# Build the application into a binary (only for use in deployments)
build:
	go build -o bin/application application.go
