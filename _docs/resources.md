# Project resources

A list of useful references used while developing this API.

## Vendors

### Docker

* [Docker Hub - PostgreSQL](https://hub.docker.com/_/postgres)
* [Limiting a container's access to memory](https://docs.docker.com/config/containers/resource_constraints/#limit-a-containers-access-to-memory)
    * This covers the concepts in use with the `shm_size` field in `docker-compose.yml`

### GNU

* [`Makefile` documentation](https://www.gnu.org/software/make/manual/make.html)

### Golang

* [`defer`, `panic` and `recover`](https://go.dev/blog/defer-panic-and-recover)

### Gorm

* [Homepage](https://gorm.io/index.html)
* [Connecting to the Database - PostgreSQL](https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL)

### PostgreSQL

* [Server start](https://www.postgresql.org/docs/current/server-start.html)

## Community

### Blogs

* [List of TimeZones supported by Postgres (version 11)](https://bill.harding.blog/2020/03/21/list-of-postgres-11-time-zones/)
* [REST API tutorial using Gin and Gorm](https://blog.logrocket.com/rest-api-golang-gin-gorm/)
* [`Makefile`s for Go Developers](https://tutorialedge.net/golang/makefiles-for-go-developers/)
* [Makefile tutorial](https://makefiletutorial.com/)
* [How to create a controller](https://letsgo-framework.github.io/guides/controllers.html#how-to-create-a-controller)
* [When to use pointers in Go](https://medium.com/@meeusdylan/when-to-use-pointers-in-go-44c15fe04eac)
* [Go in production **(Soundcloud's standard practices)**](http://peter.bourgon.org/go-in-production)
* [Go thread management](https://blog.stackademic.com/deep-dive-into-go-runtime-advanced-thread-management-explained-fd5c97b4daa4)

### GitHub

#### Documentation

* [Creating PostgreSQL service containers](https://docs.github.com/en/actions/using-containerized-services/creating-postgresql-service-containers)

#### Repositories
* [`gin-gonic/gin`](https://github.com/gin-gonic/gin) - _Gin - HTTP web framework written in Go_
* [`godotenv`](https://github.com/joho/godotenv) - _Go port of Ruby's dotenv library, which loads variables from a `.env` file_
* [`gin-swagger`](https://github.com/swaggo/gin-swagger) - _gin middleware to automatically generate RESTful API documentation with Swagger 2.0_

#### Issues
* [How to close connection in V2 (`gorm`)](https://github.com/go-gorm/gorm/issues/3145)

### Reddit

* [`httptest` recorder returns the wrong status](https://www.reddit.com/r/golang/comments/10o654j/httptest_response_recorder_returns_the_wrong/)

### StackOverflow

* [Using `.env` variables in your `docker-compose.yml` file](https://stackoverflow.com/questions/29377853/how-can-i-use-environment-variables-in-docker-compose)
* [Loading a `.env` file in a `Makefile`](https://stackoverflow.com/questions/44628206/how-to-load-and-export-variables-from-an-env-file-in-makefile)
* [Create a file using a `Makefile`](https://stackoverflow.com/questions/2667789/how-to-create-a-file-using-makefile)
* [Test naming conventions in Golang](https://stackoverflow.com/questions/15148331/test-naming-conventions-in-golang)
* [Project structure recommendations for Golang Gin projects](https://stackoverflow.com/questions/57024470/folder-structure-and-package-naming-convention-for-a-rest-api-develop-in-gin-fra)
* [When to use `os.Exit()` and `os.Panic()`](https://stackoverflow.com/questions/28472922/when-to-use-os-exit-and-panic) (short answer: _not often_)
* [Using UUID in Golang/Gorm](https://stackoverflow.com/questions/36486511/how-do-you-do-uuid-in-golangs-gorm)
* [Separating unit tests and integration tests in Go](https://stackoverflow.com/questions/25965584/separating-unit-tests-and-integration-tests-in-go)
* [Mock functions in Go](https://stackoverflow.com/questions/19167970/mock-functions-in-go)
* [Golang defaults to running tests concurrently in separate Goroutines](https://stackoverflow.com/questions/24375966/does-go-test-run-unit-tests-concurrently)
* [If you can decode JWT, how are they secure?](https://stackoverflow.com/questions/27301557/if-you-can-decode-jwt-how-are-they-secure)
* [Use of defer in go](https://stackoverflow.com/questions/47607955/use-of-defer-in-go)

### YouTube

* [Improve Your Go Tests with TestMain](https://www.youtube.com/watch?v=MAdwtwHzGP4)
