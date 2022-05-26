FROM debian:bullseye

RUN apt update -y && apt install -y build-essential cmake curl libicu-dev libz-dev git unzip

ARG CBLBUILD=couchbase-mobile-tools/cblite/build

RUN git clone --recurse-submodules https://www.github.com/couchbaselabs/couchbase-mobile-tools.git

RUN mkdir -p ${CBLBUILD} && cd ${CBLBUILD} && cmake .. && cmake --build . --parallel 4

ENV PATH="${CBLBUILD}:${PATH}"

RUN cblite --help

COPY /cblite_entrypoint.sh /cblite_entrypoint.sh

CMD ["/bin/bash", "/cblite_entrypoint.sh"]
