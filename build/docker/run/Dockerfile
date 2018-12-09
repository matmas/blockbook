FROM ubuntu:bionic

RUN apt-get update && \
    apt-get install -y libsnappy1v5 libzmq5 gettext-base && \
    apt-get clean && \
    groupadd --gid 1001 --system blockbook && \
    useradd --uid 1000 --system -g blockbook --create-home blockbook && \
    mkdir /data && \
    chown blockbook:blockbook /data

VOLUME /data

COPY build/blockbook build/docker/run/docker-entrypoint.sh /usr/local/bin/
COPY build/blockchaincfg.json.template /opt
COPY static /opt/static

# Working directory needs to contain static/
WORKDIR /opt

USER blockbook

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["blockbook", "-blockchaincfg=/home/blockbook/blockchaincfg.json","-datadir=/data","-sync","-public=:8080","-logtostderr"]

EXPOSE 8080
