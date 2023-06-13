<div style="text-align: center;" align="center">
<img src="web/public/ivory.png" alt="logo" />

# Ivory
### [postgres / patroni cluster management tool]

<img src="https://img.shields.io/github/deployments/veegres/ivory/production" alt="stable version" />
<img src="https://img.shields.io/docker/v/aelsergeev/ivory/latest?label=stable" alt="stable version" />
<img src="https://img.shields.io/docker/v/aelsergeev/ivory?label=latest" alt="latest version" />
<img src="https://img.shields.io/docker/pulls/aelsergeev/ivory" alt="docker pulls" />
</div>

<br>

Ivory is an open-source project designed to simplify and visualize work with Postgres clusters.
Initially, this tool was developed to ease the life of developers who maintain Postgres.
But I hope it will help manage and troubleshoot Postgres clusters for both developers and database administrators.

Ivory can be used as a local tool in your personal computer for your needs and as a standalone tool
in a separete virtual machine for collaborative usage.

## Get started
1. Start the docker container
2. Go to http://localhost:80
3. Do the initial configuration (Ivory will guide you)
4. Add your first cluster (by providing name and instances)
5. Start monitoring :) 

You can simply start and run container from Docker Hub or GitHub Container repository

- Docker Hub `docker run -p 80:80 --restart always aelsergeev/ivory`
- GitHub Container repository `docker run -p 80:80 --restart always ghcr.io/veegres/ivory`

## Guide

### Data
All Ivory data is located inside `/opt/data` directory. Ivory has a docker volume, it means that you won't
lose it if your container are going to be rebooted. But you need to consider mount this directory to your 
local disk if you want to save the data between different containers 
`--mount type=bind,source=YOUR_LOCAL_PATH,target=/opt/data`, or you can mount volume of the 
old container to the new one by docker flag `--volumes-from`

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
- [Good for newcomers](https://github.com/veegres/ivory/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
- [Supported versions](SECURITY.md)
- [Setup Local Environment](docker/development/README.md)
- [Build Frontend](web/README.md)
- [Build Backend](service/README.md)
