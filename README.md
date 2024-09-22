<div style="text-align: center;" align="center">
   <img src="web/public/ivory.png" alt="logo" />

   # Ivory
   ### [postgres / patroni cluster management tool]

   <img src="https://img.shields.io/github/deployments/veegres/ivory/production" alt="deployment" />
   <img src="https://img.shields.io/docker/v/aelsergeev/ivory/latest?label=stable" alt="stable version" />
   <img src="https://img.shields.io/docker/v/aelsergeev/ivory?label=latest" alt="latest version" />
   <img src="https://img.shields.io/docker/pulls/aelsergeev/ivory" alt="docker pulls" />
</div>

<br>

Ivory is an open-source project designed to simplify and visualize work with Postgres clusters.
Initially, this tool was developed to ease the life of developers who maintain Postgres.
But I hope it will help manage and troubleshoot Postgres clusters for both developers and database administrators.

Ivory allows you to use it as a local tool in your personal computer or as a standalone tool
in a separate virtual machine for collaborative usage, and helps you:
- [keep all of your clusters in one place](doc/clusters.md)
- [provides UI for all main features of patroni](doc/overview.md)
- [view and edit cluster config](doc/config.md)
- [simply execute and save template requests for troubleshooting](doc/instance.md)
- [check and clean bloat](doc/bloat.md)

## Get started
1. Start the docker container
   - Docker Hub `docker run -p 80:80 --restart always aelsergeev/ivory`
   - GitHub Container repository `docker run -p 80:80 --restart always ghcr.io/veegres/ivory`
2. Go to http://localhost:80
3. Do the initial configuration (Ivory will guide you)
4. Add your first cluster (by providing name and instances)
5. Start monitoring :)

![Demo](doc/images/demo.gif)

## Q&A

### How to update to a new version?
Unfortunately, Ivory doesn't guaranty backward compatability between minor and major releases, path releases
are questionable, but usually they contain only bug fixes and some small improvements. You can always check
in the security page backward compatibility between versions. Hence, we recommend to install it from scratch
for minor and major releases. There is some plans to introduce import/export to simplify this process. To save
your data between containers you can follow instructions that are provided bellow.

### How Ivory stores the data?
All Ivory data is located inside `/opt/data` directory. Ivory has a docker volume, it means that you won't
lose it if your container are going to be rebooted. But you need to consider mount this directory to your 
local disk if you want to save the data between different containers 
`--mount type=bind,source=YOUR_LOCAL_PATH,target=/opt/data`, or you can mount volume of the 
old container to the new one by docker flag `--volumes-from`

### How to use authentication?
Ivory can work with or without authentication. It will ask you to configure it in the initial start. Right now
Ivory supports only _Basic_ authentication with general username and password (maybe in the future support
will be added for such things like ldap / sso). Usually you don't want to use authentication when you work 
with Ivory locally, but it is recommended to use it if you use it in some VMs.

### How to run Ivory under a sub path?
Ivory offers a special environment variable, `IVORY_URL_PATH`, designed for use when running the service behind 
a reverse proxy under a sub-path. It’s important to note that the path must start with a leading slash, such 
as `/ivory`. Example: `docker run -p 80:80 -env IVORY_URL_PATH=/ivory --restart always aelsergeev/ivory`

### How to run Ivory under TLS?
You need to specify two environment variables `IVORY_CERT_FILE_PATH` and `IVORY_CERT_KEY_FILE_PATH`. Because it is
docker you need to mount these files first and then provide these variables with paths. Recommended path inside 
container is `/opt/certs`. Note that Ivory changes port to 443 when you have provided both paths. 
Example: `docker run -p 443:443 --mount type=bind,source=/etc/ssl/certs,target=/opt/certs
--env IVORY_CERT_FILE_PATH=/opt/certs/YOUR_CERT_NAME.crt --env IVORY_CERT_KEY_FILE_PATH=/opt/certs/YOUR_KEY_NAME.key 
--restart always aelsergeev/ivory`

## Contribution

If you're interested in contributing to the Ivory project, consider these options:

- [Enhancements](https://github.com/veegres/ivory/issues)
- [Good for newcomers](https://github.com/veegres/ivory/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
- [Supported versions](SECURITY.md)
- [Setup Local Environment](docker/development/README.md)
- [Build Frontend](web/README.md)
- [Build Backend](service/README.md)
