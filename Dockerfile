FROM ubuntu:22.04

EXPOSE 80
WORKDIR /opt
VOLUME /opt/data

RUN apt update
RUN apt install -y curl htop
RUN apt install -y nginx
RUN apt install -y libdbi-perl libdbd-pg-perl

# move build files to container
COPY service/build /opt/service
COPY web/build /opt/web

# move nginx settings to container
COPY docker/production/nginx.conf /etc/nginx/nginx.conf

# move docker entry file to container
COPY docker/production/entrypoint.sh /usr/local/bin/
RUN chmod 777 /usr/local/bin/entrypoint.sh

# set up pgcompacttable
RUN curl --create-dir -o /opt/tools/pgcompacttable https://raw.githubusercontent.com/dataegret/pgcompacttable/master/bin/pgcompacttable
RUN chmod +x /opt/tools/pgcompacttable
RUN ln -s /opt/tools/pgcompacttable /usr/bin

ARG IVORY_VERSION_TAG=none
ARG IVORY_VERSION_COMMIT=none
ENV IVORY_VERSION_TAG=$IVORY_VERSION_TAG
ENV IVORY_VERSION_COMMIT=$IVORY_VERSION_COMMIT

ENTRYPOINT ["entrypoint.sh"]
