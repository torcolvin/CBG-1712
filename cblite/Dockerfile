FROM debian:bullseye

RUN apt update -y && apt install -y curl unzip

ARG CBLITE_ZIP=linux-x86_64.zip

RUN curl -L -O https://github.com/couchbaselabs/couchbase-mobile-tools/releases/download/cblite-3.0.0EE-alpha/$CBLITE_ZIP && unzip $CBLITE_ZIP

RUN cblite --help

COPY /cblite_entrypoint.sh /cblite_entrypoint.sh

CMD ["/bin/bash", "/cblite_entrypoint.sh"]
