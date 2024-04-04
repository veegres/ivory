FROM postgres:16.2

RUN apt update
RUN apt -y install python3-pip

RUN pip install --break-system-packages "patroni[psycopg3,raft]"

COPY patroni-config.yml /opt
COPY patroni-initdb.sh /opt
COPY certs /opt/certs
COPY certs/ca/ca.crt /opt/certs/postgres

RUN chmod -R 777 /opt

# this is needed for postgres to establish ssl
RUN chown -R postgres:postgres /opt/certs/postgres
RUN chmod -R 600 /opt/certs/postgres/postgres-cert.key

USER postgres
WORKDIR /tmp
ENTRYPOINT patroni /opt/patroni-config.yml
