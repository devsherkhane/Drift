FROM mysql:8.0

# Copy the initialization script to the special directory
# Scripts in this directory are executed on container startup
COPY init.sql /docker-entrypoint-initdb.d/

# Ensure the database name matches what we expect
ENV MYSQL_DATABASE=drift
