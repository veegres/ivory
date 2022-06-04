# Ivory (patroni cluster visualization)

Ivory is a open-source project which is designed to simplify and visualize work with postgres clusters

Ivory allows you to do such things as:
- keep all of your cluster in one place
- monitor cluster status, do reinit and switchover
- edit config of the cluster
- clean bloat

![Demo](https://github.com/veegres/ivory/blob/master/doc/images/demo.gif)

More examples you can find in `doc` package.

## Get started

You can simply start and run container from Docker Hub or GitHub Container repository

- Docker Hub `docker run -p 80:80 aelsergeev/ivory`
- GitHub Container repository `docker run -p 80:80 ghcr.io/veegres/ivory`

## Contributing

If you're interested in contributing to the Ivory project:

- Please, check the main ideas for enhancements in [CONTRIBUTING.md](https://github.com/veegres/ivory/blob/master/CONTRIBUTING.md)
- Learn how to set up your local environment in this [README.md](https://github.com/veegres/ivory/tree/master/docker/development)
- Learn how to run frontend in this [README.md](https://github.com/veegres/ivory/blob/master/web/README.md)
- Learn how to run backend in this [README.md](https://github.com/veegres/ivory/blob/master/service/README.md)
- Look through our style guide _(in progress)_
