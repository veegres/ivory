FROM ubuntu:22.04

EXPOSE 80
WORKDIR /opt
VOLUME /opt/data

RUN apt update
RUN apt install -y curl htop
RUN apt install -y libdbi-perl libdbd-pg-perl

# move build files to container
COPY service/build /opt/service
COPY web/build /opt/web

# move docker entry file to container
COPY docker/production/entrypoint.sh /usr/local/bin/
RUN chmod 777 /usr/local/bin/entrypoint.sh

# set up pgcompacttable
RUN curl --create-dir -o /opt/tools/pgcompacttable https://raw.githubusercontent.com/dataegret/pgcompacttable/master/bin/pgcompacttable
RUN ln -s /opt/tools/pgcompacttable /usr/bin

RUN chmod -R 777 /opt

ARG IVORY_URL_PATH=""
ENV IVORY_URL_PATH=$IVORY_URL_PATH

ARG IVORY_VERSION_TAG=none
ENV IVORY_VERSION_TAG=$IVORY_VERSION_TAG

ARG IVORY_VERSION_COMMIT=none
ENV IVORY_VERSION_COMMIT=$IVORY_VERSION_COMMIT

ENTRYPOINT ["entrypoint.sh"]
