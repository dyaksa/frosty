FROM postgres:latest

COPY ./migrations/*.sql /docker-entrypoint-initdb.d/
COPY ./*.sh /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/*.sh
