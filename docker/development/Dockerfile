FROM postgres:14.2

RUN apt update
RUN apt -y install python3-pip

RUN pip install "psycopg[binary]>=3.0.0"
RUN pip install "patroni[raft]"

COPY patroni-config.yml /opt/
COPY patroni-initdb.sh /opt/

USER postgres
WORKDIR /tmp
ENTRYPOINT patroni /opt/patroni-config.yml
