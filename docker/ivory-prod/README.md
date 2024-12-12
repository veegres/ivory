## Docker environment for production

You can **build** and **run** project locally by executing this commands.

1. `make build` - build binary for service (in _service_ package)
2. `yarn run build` or `npm run build` - build web project (in _web_ package)
3. `docker build -t ivory .` - create image for docker (in _root_ package)
4. `docker run -d -p 80:80 --name ivory ivory` - run docker container

Please, be aware that it won't work with **development environment**
while you do not connect this container to network of the environment. You can do that by executing this command

5. (optional) `docker network connect development_dev-patroni ivory` - connect to dev environment (`patroni[1-3]:[8001-8003]`)

Dockerfile is located in root path, because of docker restrictions

In production, it will work only with full domains name like _google.com_, just _google_ won't work cause container
doesn't know anything about your local machine network

### Environment variables

- `IVORY_URL_PATH` - can be set by user, `default: /`
- `IVORY_STATIC_FILES_PATH` - set in entrypoint.sh file
- `IVORY_VERSION_TAG` - should be set by build system
- `IVORY_VERSION_COMMIT` - should be set by build system
