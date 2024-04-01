FROM postgres:16.2

RUN apt update
RUN apt -y install python3-pip

RUN pip install --break-system-packages "patroni[psycopg3,raft]"

COPY patroni-config.yml /opt/
COPY patroni-initdb.sh /opt/
COPY certs /opt/certs/

RUN chmod -R 777 /opt

USER postgres
WORKDIR /tmp
ENTRYPOINT patroni /opt/patroni-config.yml
