FROM redis:6.2.6

COPY  ./redis.conf /usr/local/etc/redis/redis.conf

COPY ./tls/redis.crt /usr/local/etc/redis/redis.crt
COPY ./tls/redis.key /usr/local/etc/redis/redis.key
COPY tls/ca.crt /usr/local/etc/redis/ca.crt

RUN chown -R redis:redis /usr/local/etc/redis

USER root
CMD ["redis-server", "/usr/local/etc/redis/redis.conf"]
