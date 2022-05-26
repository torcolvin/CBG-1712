FROM couchbase/server:7.1.0

#RUN yum update -y && yum install -y jq
RUN apt update -y && apt install -y jq

RUN jq -h

COPY ./server-sandbox/scripts/configure-node.sh /etc/service/config-couchbase/run
RUN chown -R couchbase:couchbase /etc/service
COPY ./server-sandbox/scripts/create-index.json /opt/couchbase

HEALTHCHECK  --start-period=15s  \
	CMD couchbase-cli server-info  -u Administrator -p password -c http://localhost:8091 || exit 1
