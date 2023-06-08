# Ivory (postgres/patroni cluster management tool)

Ivory is an open-source project which is designed to simplify and visualize work with postgres clusters.
Initially, this tool was developed to simplify the life of developers who, in their work, maintain postgres, but
I hope it will help to manage and troubleshoot postgres clusters not only for developers, but for more
experienced people with postgres.

## Get started

You can simply start and run container from Docker Hub or GitHub Container repository

- Docker Hub `docker run -p 80:80 --restart always aelsergeev/ivory`
- GitHub Container repository `docker run -p 80:80 --restart always ghcr.io/veegres/ivory`

## Ivory allows you to do such things as:
- [keep all of your cluster in one place](doc/clusters.md)
- [monitor cluster status, do reinit and switchover](doc/overview.md)
- [edit cluster config](doc/config.md)
- [monitor and clean bloat](doc/bloat.md)
- [troubleshoot particular instance](doc/instance.md)

![Demo](doc/images/demo.gif)



## Contributing

If you're interested in contributing to the Ivory project:

- [Enhancements](https://github.com/veegres/ivory/issues)
- [Supported versions](SECURITY.md)
- [Setup Local Environment](docker/development/README.md)
- [Build Frontend](web/README.md)
- [Build Backend](service/README.md)
