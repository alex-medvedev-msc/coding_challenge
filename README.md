# Coding Challenge Task
## Service for funds transaction between accounts with a RESTful API.

The service allows for updating the balances of accounts and handling the following endpoint:

- API endpoint ``GET /v1/accounts`` which shows a list of accounts
- API endpoint ``GET /v1/payments`` which shows a list of payments
- API ``POST /v1/payments`` that creates a new payment with given ``from_account``, ``amount``, ``to_account``.

## Tech stack

- go 1.10+
- godep
- postgres 10.4
- docker
- docker-compose

## Design and architecture

Key problem of updating account's balance solved with SELECT ... FOR UPDATE on postgres level.

For simplicity of deployment all bundled in docker and builds also in docker container

There are only unit-tests for key functionality and integration tests which use postgres.
In this particular project functional tests are not needed because some critical parts of project functionality
done on the db level

There is no filtering either ``GET /v1/accounts`` or ``GET /v1/payments``
because there was no such requirements in task

Account balance and payment amount are uncapped, but cannot be negative

There is a check for the currency mismatch between two accounts

Additional endpoints were not implemented although it's implementation can be helpful for testing purposes

## Start

``make run``

**Requirements**:

- docker
- docker-compose

## Unit tests

``make unit_test``

**Requirements**:

- go 1.10+

## Integration tests

``make test``

**Requirements**:

- go 1.10+
- docker
- docker-compose

## Bench

``make bench``

**Requirements**:

- go 1.10+
- docker
- docker-compose

## Lint

``make lint``

**Requirements**:

- go 1.10+
- gometalinter.v2 with installed linters

