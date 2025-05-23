FROM mysql:latest
COPY init.sql /docker-entrypoint-initdb.d/

ENV MYSQL_ROOT_PASSWORD=arti2002
ENV MYSQL_DATABASE=forum_db
ENV MYSQL_HOST=mysql