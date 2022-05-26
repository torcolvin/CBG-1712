#!/bin/bash

set -eux -o pipefail

./entrypoint.sh couchbase-server &

CURL="curl -u Administrator:password --fail"
CURL_RETRY="$CURL --retry-all-errors --connect-timeout 5 --max-time 2 --retry 20 --retry-delay 0 --retry-max-time 120"
SERVER_ADDR=http://localhost:8091
COUCHBASE_CLI="couchbase-cli"
COUCHBASE_CLI_OPTS="--username Administrator --password password --cluster $SERVER_ADDR"

$CURL_RETRY $SERVER_ADDR/pools/default/buckets

$CURL http://127.0.0.1:8091/nodes/self/controller/settings -d 'path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&' -d 'index_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&' -d 'cbas_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&' -d 'eventing_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata&'
$CURL http://127.0.0.1:8091/node/controller/setupServices -d 'services=kv%2Cn1ql%2Cindex'
$CURL http://127.0.0.1:8091/pools/default -d 'memoryQuota=3072' -d 'indexMemoryQuota=3072' -d 'ftsMemoryQuota=256'
$CURL http://127.0.0.1:8091/settings/web -d 'password=password&username=Administrator&port=SAME'
$CURL http://localhost:8091/settings/indexes -d indexerThreads=4 -d logLevel=verbose -d maxRollbackPoints=10 \
    -d storageMode=plasma -d memorySnapshotInterval=150 -d stableSnapshotInterval=40000

$COUCHBASE_CLI bucket-create $COUCHBASE_CLI_OPTS --bucket testBucket --bucket-type couchbase --bucket-ramsize 512
