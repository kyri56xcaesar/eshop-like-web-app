FROM postgres:latest
ENV POSTGRES_DB orders_db
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD pass123
EXPOSE 5433:5432
COPY ./init_products_table.sql /docker-entrypoint-initdb.d/