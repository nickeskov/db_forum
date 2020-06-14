# Step 1. build_step
FROM golang:1.14-stretch AS build_step
WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o my_db_forum cmd/main.go

# Step 2. release_step
FROM ubuntu:18.04 AS release_step

MAINTAINER Nicholas Eskov

ENV DEBIAN_FRONTEND=noninteractive

ENV PGVER 12

RUN apt -y update && \
    apt install -y wget gnupg && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    echo "deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main" >> /etc/apt/sources.list.d/pgdg.list && \
    apt -y update

RUN apt -y update && apt install -y \
        postgresql-$PGVER \
    && rm -rf /var/lib/apt/lists/*

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Copy database.sql script and postresql.conf custom config
COPY --from=build_step /app/configs/database/sql/database.sql /app/database.sql
COPY --from=build_step /app/configs/database/sql/postgresql.conf /app/postgresql.conf

# Create a PostgreSQL role named ``my_db_forum`` with ``my_db_forum`` as the password and
# then create a database `my_db_forum` owned by the ``my_db_forum`` role.
RUN service postgresql start && \
    psql --command "CREATE USER my_db_forum WITH SUPERUSER PASSWORD 'my_db_forum';" && \
    createdb -O my_db_forum my_db_forum && \
    psql my_db_forum --echo-all --file /app/database.sql && \
    service postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
# Add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
# Add our posgres configuration settings
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf && \
    echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf && \
    cat /app/postgresql.conf >> /etc/postgresql/$PGVER/main/conf.d/custom_postgresql.conf

# Back to the root user
USER root

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Expose the PostgreSQL port
EXPOSE 5432

# Expose server port
EXPOSE 5000

COPY --from=build_step /app/my_db_forum /app/

WORKDIR /app

CMD service postgresql start && ./my_db_forum
