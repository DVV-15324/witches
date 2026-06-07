FROM mcr.microsoft.com/mssql/server:2022-latest

COPY entrypoint.sh /entrypoint.sh

USER root

RUN chmod +x /entrypoint.sh

USER mssql

EXPOSE 1433

CMD ["/bin/bash", "/entrypoint.sh"]
