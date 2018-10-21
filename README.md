# fstrie-web
A basic web app example of using fstrie as an in-memory database and web app interface.

The model is a proof of concept for using fstrie to back a web app.

There are two kinds files, virtual, and real. Both fulfill a shared interface and are easily extensible to have more complex functionality.

New real files will be added to the db as they come into existence. There is not presently a routine to remove files from the db, though this would be trivial to implement.

## Dependencies

- [fstrie](https://github.com/henesy/fstrie)

## Build

	go build

## Run

	./fstrie-web

## Usage

Navigate to [localhost:8080](http://localhost:8080).

Default virtual page is at [root](http://localhost:8080/).

Default real page is at [default.html](http://localhost:8080/default.html).

Default 404 page is everywhere else!
