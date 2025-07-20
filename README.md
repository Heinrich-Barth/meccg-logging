# Game monitor

This program fetches a list of active games every 5mins and saves them to disk.

The games will be saved by their date and `ìd` in json format.

You can query all saved game information via

`http://localhost:{port}/games`

and a specific game via

`http://localhost:{port}/games/game-file-name.json`

## Start

To start the app, run 

```bash
go run .
```

## Build and execute

Use these commands

```bash
go build -o dist/app .
./dist/app
```
