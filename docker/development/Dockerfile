FROM postgres:14.2

RUN apt update
RUN apt -y install python3-pip

RUN pip install "psycopg[binary]>=3.0.0"
RUN pip install "patroni[raft]"

COPY patroni-config.yml /opt/
COPY patroni-initdb.sh /opt/
COPY certs /opt/certs/

RUN chmod -R 777 /opt

USER postgres
WORKDIR /tmp
ENTRYPOINT patroni /opt/patroni-config.yml
