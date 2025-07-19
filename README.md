# Game monitor

This program fetches a list of active games every 5mins and saves them to disk.

The games will be saved by their `ìd` in json format.

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
