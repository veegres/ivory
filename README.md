# Ivory (postgres / patroni cluster management tool)

Ivory is an open-source project which is designed to simplify and visualize work with postgres clusters.
Initially, this tool was developed to simplify the life of developers who, in their work, maintain postgres, but
I hope it will help to manage and troubleshoot postgres clusters not only for developers, but for more
experienced people with postgres.

## Get started

You can simply start and run container from Docker Hub or GitHub Container repository

- Docker Hub `docker run -p 80:80 --restart always aelsergeev/ivory`
- GitHub Container repository `docker run -p 80:80 --restart always ghcr.io/veegres/ivory`

## Guide

### Authentication
Ivory can work with or without authentication. It can be configured in the initial start. Right now
Ivory supports only _Basic_ authentication with general username and password (maybe in the future support
will be added for such things like ldap / sso)
Usually you don't want to use authentication when you work with Ivory locally, but it is recommended
to use it if you put it globally.

### Usage
- [keep all of your cluster in one place](doc/clusters.md)
- [monitor cluster status, do reinit and switchover](doc/overview.md)
- [edit cluster config](doc/config.md)
- [monitor and clean bloat](doc/bloat.md)
- [troubleshoot particular instance](doc/instance.md)

![Demo](doc/images/demo.gif)


## Contribution

If you're interested in contributing to the Ivory project:

- [Enhancements](https://github.com/veegres/ivory/issues)
- [Supported versions](SECURITY.md)
- [Setup Local Environment](docker/development/README.md)
- [Build Frontend](web/README.md)
- [Build Backend](service/README.md)
