include .env
export

postgres:
	sudo docker run --name postgres -p 5432:5432 -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d postgres

createdb: 
	sudo docker exec -it postgres createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	sudo docker exec -it postgres dropdb  $(DB_NAME)

start:
	docker-compose up

stop:
	docker-compose down

.PHONY: postgres createdb dropdb start stop
