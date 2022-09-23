# Ivory (patroni cluster management)

Ivory is an open-source project which is designed to simplify and visualize work with postgres clusters

Ivory allows you to do such things as:
- [keep all of your cluster in one place](https://github.com/veegres/ivory/blob/master/doc/clusters.md)
- [monitor cluster status, do reinit and switchover](https://github.com/veegres/ivory/blob/master/doc/overview.md)
- [edit config of the cluster](https://github.com/veegres/ivory/blob/master/doc/config.md)
- [clean bloat](https://github.com/veegres/ivory/blob/master/doc/bloat.md)

![Demo](https://github.com/veegres/ivory/blob/master/doc/images/demo.gif)

## Get started

You can simply start and run container from Docker Hub or GitHub Container repository

- Docker Hub `docker run -p 80:80 --restart always aelsergeev/ivory`
- GitHub Container repository `docker run -p 80:80 --restart always ghcr.io/veegres/ivory`

### Docker environment variables

You can customize Ivory by providing some environment variables to the docker container

```
IVORY_COMPANY_LABEL  - set up custom company name at the left top corner
IVORY_AUTHENTICATION - set up authentication type "none" or "basic"
IVORY_BASIC_USER     - set up user for basic authentication field is mandatory
IVORY_BASIC_PASSWORD - set up password for basic authentication field is mandatory
```

## Contributing

If you're interested in contributing to the Ivory project:

- Please, check the main ideas for enhancements in [TODO.md](https://github.com/veegres/ivory/blob/master/TODO.md)
- Learn how to set up your local environment in this [README.md](https://github.com/veegres/ivory/tree/master/docker/development)
- Learn how to run frontend in this [README.md](https://github.com/veegres/ivory/blob/master/web/README.md)
- Learn how to run backend in this [README.md](https://github.com/veegres/ivory/blob/master/service/README.md)
- Look through our style guide _(coming soon)_
