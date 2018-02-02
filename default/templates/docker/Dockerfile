FROM ubuntu

COPY ./settings.toml /etc/{{.ServiceName}}/settings.toml
COPY ./{{.ServiceName}} /usr/bin/{{.ServiceName}}
COPY ./entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 10000
EXPOSE 10001
EXPOSE 8080
