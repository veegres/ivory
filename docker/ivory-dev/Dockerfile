FROM postgres:16.2

RUN apt update
RUN apt -y install python3-pip

RUN pip install --break-system-packages "patroni[psycopg3,raft]"

COPY patroni-config.yml /opt
COPY patroni-initdb.sh /opt
# it'll copy and client certs to the service, we don't care, because it is for test only
COPY certs /opt/certs

# this is needed for postgres and patroni to get access to the files
RUN chown -R postgres:postgres /opt/certs
# this is needed for postgres to read the key, the view should be resticted by postgres rules
RUN chmod -R 600 /opt/certs/server/server.key

USER postgres
# unclear why it doesn't want to work without making `tmp` as working dear
WORKDIR /tmp
ENTRYPOINT patroni /opt/patroni-config.yml
