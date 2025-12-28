<div style="text-align: center;" align="center">
   <img src="web/public/ivory.png" alt="logo" />

   # Ivory
   ### [postgres / patroni cluster management tool]

   <img src="https://img.shields.io/github/deployments/veegres/ivory/production?style=flat-square&link=https%3A%2F%2Fgithub.com%2Fveegres%2Fivory%2Fdeployments%2Fproduction" alt="deployment" />
   <img src="https://img.shields.io/docker/v/veegres/ivory/latest?label=stable&style=flat-square&link=https%3A%2F%2Fhub.docker.com%2Fr%2Fveegres%2Fivory" alt="stable version" />
   <img src="https://img.shields.io/docker/v/veegres/ivory?label=latest&style=flat-square&link=https%3A%2F%2Fhub.docker.com%2Fr%2Fveegres%2Fivory" alt="latest version" />
   <img src="https://img.shields.io/docker/pulls/veegres/ivory?style=flat-square&link=https%3A%2F%2Fhub.docker.com%2Fr%2Fveegres%2Fivory" alt="docker pulls" />
</div>

<br>

Ivory is an open-source project designed to simplify and visualize work with Postgres clusters.
Initially, this tool was developed to ease the life of developers who maintain Postgres.
But I hope it will help manage and troubleshoot Postgres clusters for both developers and database administrators.

Ivory allows you to use it as a local tool in your personal computer or as a standalone tool
in a separate virtual machine for collaborative usage and helps you:
- [keep all of your clusters in one place](doc/clusters.md)
- [provides UI for all main features of patroni](doc/overview.md)
- [view and edit cluster config](doc/config.md)
- [simply execute and save template requests for troubleshooting](doc/instance.md)
- [check and clean bloat](doc/bloat.md)

<div align="center">
  <h2>🌟 Support This Project! 🌟</h2>
</div>

If you found this project helpful, interesting, or inspiring, please consider giving it a **star** ⭐! Your support helps:

✅ **Increase visibility** – More people can discover and benefit from this project.  
✅ **Boost motivation** – It encourages us to keep improving and adding new features.  
✅ **Show appreciation** – A small gesture that means a lot to open-source creators!

Thank you for being part of this journey! 🚀

> *"Alone we can do so little; together we can do so much."* — Helen Keller

## Get started
1. Start the docker container
   - Docker Hub `docker run -p 80:80 --restart always veegres/ivory`
   - GitHub Container repository `docker run -p 80:80 --restart always ghcr.io/veegres/ivory`
2. Go to http://localhost:80
3. Do the initial configuration (Ivory will guide you)
4. Add your first cluster (by providing name and instances)
5. Start monitoring :)

![Demo](doc/images/demo.gif)

## Q&A

### How to update to a new version?
Ivory now provides Backup/Restore functionality for migrating your data between versions. You can backup your
configuration (clusters, queries, permissions) from the Settings page and restore it in a new version. The
backup format is backward compatible and designed to work across different versions.

Alternatively, you can mount the data directory between containers (see instructions below), though this
approach typically works only for patch releases. For minor and major version updates, always check the
[backward compatibility page](SECURITY.md) and prefer using the Backup/Restore feature for safer migration.

### How Ivory stores the data?
All Ivory data is located inside `/opt/ivory/data` directory. Ivory has a docker volume, it means that you won't
lose it if your container is going to be rebooted. But you need to consider mounting this directory to your 
local disk if you want to save the data between different containers 
`--mount type=bind,source=YOUR_LOCAL_PATH,target=/opt/ivory/data`, or you can mount volume of the 
old container to the new one by docker flag `--volumes-from`

### How to use authentication?
Ivory can work with or without authentication. It will ask you to configure it during the initial setup.
Ivory supports multiple authentication methods:
- **Basic** - Simple username and password authentication
- **LDAP** - Integration with LDAP directories
- **OIDC/SSO** - Single Sign-On via OpenID Connect

You can safely provide your secrets to Ivory for SSO and LDAP configuration. Ivory encrypts all secrets
using your secret word. Therefore, make sure to select the appropriate application configuration in your SSO provider.
As well, Ivory requires the _profile_ or _email_ scopes from the SSO provider in order to retrieve user information.

Additionally, Ivory includes a comprehensive permission system that allows you to control access at a granular
level (view/manage clusters, execute queries, manage configurations, etc.). You can manage user permissions
from the Settings page after authentication is enabled.

Usually you don't want to use authentication when working with Ivory locally, but it is recommended when
deploying it in shared environments or VMs.

### How to run Ivory under a sub path?
Ivory offers a special environment variable, `IVORY_URL_PATH`, designed for use when running the service behind
a reverse proxy under a sub-path. It's important to note that the path must start with a leading slash, such
as `/ivory`. Example: `docker run -p 80:80 -env IVORY_URL_PATH=/ivory --restart always veegres/ivory`

### How to run Ivory under TLS?
You need to specify two environment variables `IVORY_CERT_FILE_PATH` and `IVORY_CERT_KEY_FILE_PATH`. Because it is
a docker environment, you need to mount these files first and then provide these variables with paths. Recommended path inside
container is `/opt/certs`. Note that Ivory changes port to 443 when you have provided both paths.
Example: `docker run -p 443:443 --mount type=bind,source=YOUR_CERTS_PATH,target=/opt/certs
--env IVORY_CERT_FILE_PATH=/opt/certs/YOUR_CERT_NAME.crt --env IVORY_CERT_KEY_FILE_PATH=/opt/certs/YOUR_KEY_NAME.key
--restart always veegres/ivory`

## Contribution

If you're interested in contributing to the Ivory project, consider these options:

- [Enhancements](https://github.com/veegres/ivory/issues)
- [Good for newcomers](https://github.com/veegres/ivory/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
- [Supported versions](SECURITY.md)
- [Setup Local Environment](docker/ivory-dev/README.md)
- [Build Frontend](web/README.md)
- [Build Backend](service/README.md)
