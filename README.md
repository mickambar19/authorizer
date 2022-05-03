# Authorizer

[![Coverage](https://img.shields.io/badge/coverage-100%25-green)](https://img.shields.io/badge/coverage-100%25-green)

A cli application written in Go(Golang) that reads events and authorizes transactions for an account following a set of predefined rules.

## Contents

- [Authorizer](#authorizer)
  - [Contents](#contents)
  - [Getting started](#getting-started)
  - [Documentation](#documentation)
  - [Unit tests](#unit-tests)
  - [e2e tests](#e2e-tests)

## Getting started

1. You first need [Go](https://golang.org/) installed (**version 1.14+ is required**)

2. Install the necessary dependencies for the authorizer by executing the following command:

```sh
  $ go mod tidy
```

4. Create a file with the lines that the authorizers has to process

```
  You can create your own or use the already created ones that are located under e2e/cases/\*/input
```

3. Run the project and pass the file with the events you wan the authorizer to process

```sh
  $ go run src/main.go < e2e/cases/account-initialize/input
```

In case you want to compile the authorizer, you can run the following command:

```sh
  $ go build -o authorizer src/main.go
```

Let's pass an input file to our compiled version

```sh
  $ ./authorizer < e2e/cases/account-initialize/input
```

## Documentation

The Authorizer is created using clean architecture and implementing the mediator patter in this case represented by the controller.

### Why Clean architecture?

This way we can make the code more flexible for new business rules without creating a side effect in any other artifact.

- src
  - controller
    - Retrieves and parses data
    - Coordinates what to do with the input data
  - usecase
    - Business rules that should be applied
  - service
    - Methods that will help to manage the model
  - model
    - Representation of the objects that are going to be used

## Unit tests

Usecases are the only parts that were fully tested as it is the part that contains the business rules.

Libraries used: Gomock & testify

Gomock was used to create the corresponding mocks for the interfaces used

Testify was used for more readable assertions

1. Generate the mock files that will help us out with the tests

```sh
  $ make generate-mock
```

2. Run all the tests.

```sh
  $ go test ./...
  # In case you want it in verbose mode
  $ go test ./... -v
```

We can create the html coverage in case you want to check deeply which lines were covered by the tests

```sh
  $ make create-test-html-coverage
```

## e2e Tests

The e2e tests are created using BATS(Bash Automated Testing System) since the only assertion that we need is to compare the stdout from the authorizer with the corresponding output.

Prerequisites

- node version 14.17.0+

  1.Move to the e2e directory

```sh
  $ cd e2e
```

2.Install the dependencies

```sh
  $ npm i
```

3.Run the unit tests

```sh
  $ npm run test
```

### How to add a new e2e test

Folder structure:

- e2e
  - cases
    - CaseFolder
      - input
      - output
  - tests
    - authorizer.bats

1. You have to create a new folder under **cases** with the corresponding new test name using snakecase.

2. Generate the input and output file with the content you expect to send and retrieve from the authorizer accordingly.

3. Add the new test case to the bats file

```sh
@test "Case Name" {
  helper "case-name" # The name of the folder created in step 1
}
```
