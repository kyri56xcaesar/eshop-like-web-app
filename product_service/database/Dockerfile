FROM postgres:latest
ENV POSTGRES_DB products_db
ENV POSTGRES_PASSWORD pass123
COPY ./init_products_table.sql /docker-entrypoint-initdb.d/
