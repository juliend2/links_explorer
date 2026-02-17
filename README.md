Links crawler for Status Machine
================================

## Building

    make build

It uses GO15VENDOREXPERIMENT=1 , so we can use the vendored
versions of the libraries we use. This makes our builds stable, and not subject
to breakage if a library updates in a non-backward-compatible way.

## Testing

Running unit tests:

    go test

Testing that the crawler does work as expected (integration testing):

    make test

## Dependency management

We use [gvt](https://github.com/FiloSottile/gvt) for dependency management.

    go get -u github.com/FiloSottile/gvt

Now, to install dependencies, use this:

    gvt fetch github.com/PuerkitoBio/gocrawl

## Go version

Since we're using GO15VENDOREXPERIMENT, this projects necessitate go 1.5.+ or
higher (assuming that the GO15VENDOREXPERIMENT behaviour will still work after
1.5).

## Data that it receives

url: https://www.juliendesrosiers.com/
site_id: 2
snapshot_id: 1

## TODO

See: https://trello.com/c/mKT4NjcH/296-rewrite-the-links-crawler-in-go

