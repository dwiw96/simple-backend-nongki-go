# pgRun:
# 	sudo docker compose -f ./docker-compose.yml up -d
pgExec:
	sudo docker exec -it pg_nongki_container psql -U dwiw nongki_db
pgStop:
	sudo docker container stop pg_nongki_container

pgMigrateUsers:
	sudo docker exec -i pg_nongki_container psql -U dwiw nongki_db < internal/utils/postgres/users.sql
pgMigrateDrop:
	sudo docker exec -i pg_nongki_container psql -U dwiw nongki_db < internal/utils/postgres/drop_all_tables.sql
pgMigratePrivKey:
	sudo docker exec -i pg_nongki_container psql -U dwiw nongki_db < internal/utils/postgres/priv_key.sql