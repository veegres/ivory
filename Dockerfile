FROM ubuntu:24.04

# keep apt command as one-liner to avoid caching issues
RUN apt update && apt install -y --no-install-recommends \
    ca-certificates \
    curl htop \
    libdbi-perl libdbd-pg-perl

EXPOSE 80
EXPOSE 443
WORKDIR /opt/ivory
VOLUME /opt/ivory/data

ARG IVORY_URL_PATH=""
ENV IVORY_URL_PATH=$IVORY_URL_PATH

ARG IVORY_VERSION_TAG=none
ENV IVORY_VERSION_TAG=$IVORY_VERSION_TAG

ARG IVORY_VERSION_COMMIT=none
ENV IVORY_VERSION_COMMIT=$IVORY_VERSION_COMMIT

# move build files to container
COPY web/build /opt/ivory/web
COPY service/build /opt/ivory/service
RUN chmod +x /opt/ivory/service/ivory

# move docker entry file to container
COPY docker/ivory-prod/entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/entrypoint.sh

# set up pgcompacttable
RUN curl --create-dirs -o /opt/ivory/tools/pgcompacttable https://raw.githubusercontent.com/dataegret/pgcompacttable/master/bin/pgcompacttable
RUN chmod +x /opt/ivory/tools/pgcompacttable
RUN ln -s /opt/ivory/tools/pgcompacttable /usr/bin

ENTRYPOINT ["entrypoint.sh"]
