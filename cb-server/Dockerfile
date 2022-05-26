FROM couchbase/server:7.1.0-aarch64

RUN ./entrypoint.sh couchbase-server &

COPY ./bootstrap_db.sh /bootstrap_db.sh

RUN ["/bin/bash", "/bootstrap_db.sh"]

