# Copyright 2018-Present Couchbase, Inc.
#
# Use of this software is governed by the Business Source License included in
# the file licenses/BSL-Couchbase.txt.  As of the Change Date specified in that
# file, in accordance with the Business Source License, use of this software
# will be governed by the Apache License, Version 2.0, included in the file
# licenses/APL2.txt.

# Stage to build Sync Gateway binary
FROM golang:1.18

RUN apt update -y && apt install -y file
COPY ./  /CBG-1712/
RUN cd /CBG-1712/sync_gateway && go build .
RUN cd /CBG-1712 && go build ./...

ENTRYPOINT ["/CBG-1712/CBG-1712_entrypoint"]
