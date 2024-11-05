up:
	docker compose up -d

down:
	docker compose down --volumes --remove-orphans

cli:
	docker compose exec -it pg psql -U postgres -d postgres

# install postgis extension
install-postgis:
	docker compose exec -it pg sh -c "psql -d postgres -U postgres -f /postgis/init.sql"
