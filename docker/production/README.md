## Docker environment for production

You can **build** and **run** project locally by executing this commands.

1. `GOOS=linux GOARCH=amd64 go build -o build` - build binary for service (in _service_ package)
2. `yarn run build` or `npm run build` - build web project (in _web_ package)
3. `docker build -t ivory .` - create image for docker (in _root_ package)
4. `docker run -d -p 80:80 --name ivory ivory` - run docker container

Please, be aware that it won't work with **development environment**
while you do not connect this container to network of the environment. You can do that by executing this command

5. (optional) `docker network connect development_dev-patroni ivory` - connect to dev environment

Dockerfile is located in root path, because of docker restrictions

In production, it will work only with full domains name like _google.com_, just _google_ won't work cause container doesn't know anything about your local
machine network
