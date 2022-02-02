FROM ubuntu

RUN apt update
RUN apt install -y nginx curl htop

# move build files to container
COPY service/build /opt/service
COPY web/build /opt/web

# move nginx settings to container
COPY docker/production/nginx.conf /etc/nginx/nginx.conf

# move docker entry file to container
COPY docker/production/entrypoint.sh /usr/local/bin/
RUN chmod 777 /usr/local/bin/entrypoint.sh

EXPOSE 80
WORKDIR /opt
ENTRYPOINT ["entrypoint.sh"]
