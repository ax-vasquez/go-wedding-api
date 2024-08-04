# Testing

This repository aims for 95% code coverage and uses [CodeCov](https://about.codecov.io/) to measure coverage.

## Integration Tests

The majority of this repository's logic is covered using integration tests. Integration test files are named in the format
`*_integration_test.go` and only run with the [build contstraint](https://pkg.go.dev/go/build#hdr-Build_Constraints) ("tag") `integration`.

### Rationale

**Why emphasize integration tests over unit tests?**
* Code in `./controllers` and `./models` doesn't employ much custom logic, but rather straightforward `gin` and `gorm` implementations
* In the opinion of the lead developer for this repo, testing working logic is best done with "real-ish" data to emulate a real-world environment as closely as possible
* Testing things at scale is easier when using a real database connection; no need to mock (and, by extension, maintain) unwieldy data responses

### Running integration tests

To run the integration tests locally, run the following command from the root of the project repo:

```sh
go test -tags=integration ./... -v
```
* Runs all `*_test.go` files with the `integration` build tag
