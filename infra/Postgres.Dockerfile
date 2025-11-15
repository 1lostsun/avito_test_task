FROM postgres:14-alpine

RUN apk add --no-cache curl jq

COPY docker-entrypoint-vault.sh /docker-entrypoint-vault.sh
RUN chmod +x /docker-entrypoint-vault.sh

ENTRYPOINT ["docker-entrypoint-vault.sh"]
CMD ["postgres"]