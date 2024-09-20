# Backend development

## Build

Use `go mod tidy` to download all libraries.

In the project directory, you can run:

- `go build ivory` - Runs the app in the development mode

## Environment variables

- `IVORY_URL_PATH` - change sub path under which Ivory serve its routes, `default: /`
- `IVORY_STATIC_FILES_PATH` - add static files to host, `default: empty` 
- `IVORY_VERSION_TAG` - change version tag, `default: v0.0.0`
- `IVORY_VERSION_COMMIT` - change version commit, `default: 000000000000`

## Learn More

This app is build by:

- [Bolt](https://github.com/boltdb/bolt) - Embedded key/value database
- [Gin](https://github.com/gin-gonic/gin) - Web framework

