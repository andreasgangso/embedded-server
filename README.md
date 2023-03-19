# Host a static directory locally

### Description

Builds a single binary that has the static/ dir embedded and will serve it when you run the app.
Configurable http and https port on boot, and opens a browser.
Https certs are self-signed, as this is meant as a easy local serve thing
Half of this was done by chatgpt-4.

### Usefulness

In hindsight I don't really see this too useful because in case something goes wrong
you probably want to send the static files alongside this to the recipient anyway.
So it's better to just have the files in a folder next to it..

You can also use something like [Tauri](https://github.com/tauri-apps/tauri) if the goal is to
host a webapp, as it's actually very easy to set up.

### Usage

- Place something in static/
- Run `make install` to install [go.rice](https://github.com/GeertJohan/go.rice)
- Run `make build` to build into bin/
