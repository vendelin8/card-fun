## card-fun
card-fun is a demo about creating decks of cards, and drawing items from it.

The project uses `task` to make your life easier. If you're not familiar with Taskfile, you can take a look at [this quickstart guide](https://taskfile.dev/).

## Overview

### Create a new Deck

PUT request to /v1/create

It can be done by either:
* creating a deck from a given list of cards like `cards=AS,KD,AC,2C,KH` as query parameter.
* creating a full deck by leaving it out

Then it may be shuffled on request by adding `shuffled=true` as query parameter creating a full deck.

It returns a JSON response like the following with a random UUID as a deck identifier and number of cards.

```
{
    "deck_id": "a251071b-662f-44b6-ba11-e24863039c59",
    "shuffled": false,
    "remaining": 30
}
```

It may return errors with status `400` if there are duplicated cards, of some cards are not found or in a wrong format.

### Open a Deck

GET request to /v1/open

There is a mandatory parameter in the query for deck identifier in `deck_id`.

The response looks like the following as JSON.

```
{
    "deck_id": "a251071b-662f-44b6-ba11-e24863039c59",
    "shuffled": false,
    "remaining": 3,
    "cards": [
        {
            "value": "ACE",
            "suit": "SPADES",
            "code": "AS"
        },
        ...
    ]
}
```

It may result in an error if there's no such deck.

### Draw card(s)

POST request to /v1/draw

There is a mandatory parameter in the query for deck identifier in `deck_id`.
Another optional parameter in query is `count` that may be any positive integer until the size of the deck. If not set, it's default is 1.

The response looks like the following with the drawn list of cards:

```
{
    "cards": [
        {
            "value": "QUEEN",
            "suit": "HEARTS",
            "code": "QH"
        },
        {
            "value": "4",
            "suit": "DIAMONDS",
            "code": "4D"
        }
    ]
}
```

It may result in an error if there's no such deck, or if the deck has less cards than the requested number.

## Dependencies and install

[Taskfile](https://taskfile.dev/).
[Docker](https://docs.docker.com/get-started/) for the following commands: `task docker-build`, `task docker-run`, `task dev`, `task docker-test`

## Usage

In order to run it, you need a running redis server. The easiest way to run it is calling the following with running docker service:
```bash
task docker-run
```

Or if you have an installed Go ecosystem with GOPATH, the following will reuse package cache, so it will be much faster:
```bash
task dev
```

Or if you have a running redis on your localhost, you may run it without docker:
```bash
task run
```

**Then open <http://localhost:3000/swagger> in your browser.**

## Test & lint

Lint:

```bash
task lint
```

Run tests with running docker service:

```bash
task docker-test
```

Run tests with running redis on localhost:

```bash
task test
```

Check test coverage in browser:

```bash
task cover
```

## REST

Swagger docs can be regenerated using `task swagger`, and can be found in `cmd/server/docs`. If you already ran it, it will be symlinked to `out/docs`.