GO=${GOROOT}/bin/go

mssql-docker:
	docker run -e 'ACCEPT_EULA=Y' -e 'MSSQL_SA_PASSWORD=admini12!' -e 'MSSQL_PID=Developer' -p 1433:1433 -d --name tg-mssql mcr.microsoft.com/mssql/server:2022-latest
