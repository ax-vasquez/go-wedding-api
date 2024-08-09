# README

The API used for our wedding site, built in Go. This is a simple API used to save wedding guest preferences
and responses using a stack built from scratch.

[![codecov](https://codecov.io/gh/ax-vasquez/wedding-site-api/graph/badge.svg?token=UH5YZFRM35)](https://codecov.io/gh/ax-vasquez/wedding-site-api)

## API Documentation

This project leverages [`gin-swagger`](https://github.com/swaggo/gin-swagger) for its API documentation.

To view the API documentation locally, start the application (`go run main.go`) and navigate to `localhost://8080/swagger/index.html`

## Local development

**Quick steps**
1. `setup.sh`
1. `docker-compose up -d`
1. `go run main.go`

### Generating the `.env` file

The provided `setup.sh` script will generate the `.env` file for you and preload defaults. You can edit the generated variable
values as-needed, as well as add to the `.env` file if you need. If you somehow get into a state where you're missing a necessary
environment variable, you can re-run `setup.sh` to add the missing variables back to your `.env` file non-destructively. It will 
not overwrite values for required variables that you may have edited.

### Starting the containers

To start the supporting docker containers, run `docker-compose up -d` from the root of the repository.

### Resetting your local database

If you get your database into a weird state, it's often simplest to just delete the database and re-create it. _To be clear, this is only a valid
option if you're running locally._ Ideally, we should be hashing all of these weird things out before we are anywhere near deploying things.

To reset the DB, you'll need to do the following steps (in this order):
1. Run `docker-compose down`
1. Delete the `./local-data` directory
1. Run `docker-compose up -d`

You need to make sure the Postgres container has stopped before deleting the `/local-data` directory because it writes to this directory while in
operation. Then, once the container is stopped (and no longer writing to the `/local-data` directory), delete the `/local-data` directory. If
you don't delete the `/local-data` directory before restarting, the Postgres container will simply re-use what's defined there and it won't be reset.

Additionally, if you delete the application database manually (currently `"gorm"`) using a `DELETE DATABASE` command, _the application
will not re-create it on the next startup and it will fail to connect_. **In cases like this, it's best to simply reset your database using the
steps above.**.
