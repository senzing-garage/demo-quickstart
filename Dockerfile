# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_GO_BUILDER=golang:1.21.4-bullseye
ARG IMAGE_FINAL=senzing/senzingapi-runtime:3.8.0

# -----------------------------------------------------------------------------
# Stage: senzingapi_runtime
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} as senzingapi_runtime

# -----------------------------------------------------------------------------
# Stage: go_builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_GO_BUILDER} as go_builder
ENV REFRESHED_AT=2024-01-30
LABEL Name="senzing/demo-quickstart-builder" \
      Maintainer="support@senzing.com" \
      Version="0.0.1"

# Copy local files from the Git repository.

COPY ./rootfs /
COPY . ${GOPATH}/src/demo-quickstart

# Copy files from prior stage.

COPY --from=senzingapi_runtime  "/opt/senzing/g2/lib/"   "/opt/senzing/g2/lib/"
COPY --from=senzingapi_runtime  "/opt/senzing/g2/sdk/c/" "/opt/senzing/g2/sdk/c/"

# Set path to Senzing libs.

ENV LD_LIBRARY_PATH=/opt/senzing/g2/lib/

# Build go program.

WORKDIR ${GOPATH}/src/demo-quickstart
RUN make build

# Copy binaries to /output.

RUN mkdir -p /output \
 && cp -R ${GOPATH}/src/demo-quickstart/target/*  /output/

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} as final
ENV REFRESHED_AT=2023-10-03
LABEL Name="senzing/demo-quickstart" \
      Maintainer="support@senzing.com" \
      Version="0.0.1"

# Copy files from prior stage.

COPY ./rootfs /
COPY --from=go_builder "/output/linux-amd64/demo-quickstart" "/app/demo-quickstart"

# Install packages via apt.

RUN apt update \
 && apt -y install \
      gnupg2 \
      jq \
      libodbc1 \
      postgresql-client \
      unixodbc \
 && rm -rf /var/lib/apt/lists/*

# Install Java-11.

RUN mkdir -p /etc/apt/keyrings \
 && wget -O - https://packages.adoptium.net/artifactory/api/gpg/key/public > /etc/apt/keyrings/adoptium.asc

RUN echo "deb [signed-by=/etc/apt/keyrings/adoptium.asc] https://packages.adoptium.net/artifactory/deb $(awk -F= '/^VERSION_CODENAME/{print$2}' /etc/os-release) main" >> /etc/apt/sources.list

RUN export STAT_TMP=$(stat --format=%a /tmp) \
 && chmod 777 /tmp \
 && apt update \
 && apt -y install temurin-11-jdk \
 && chmod ${STAT_TMP} /tmp \
 && rm -rf /var/lib/apt/lists/*

# Runtime environment variables.

ENV LD_LIBRARY_PATH=/opt/senzing/g2/lib/
ENV SENZING_API_SERVER_ALLOWED_ORIGINS='*'
ENV SENZING_API_SERVER_BIND_ADDR='all'
ENV SENZING_API_SERVER_ENABLE_ADMIN='true'
ENV SENZING_API_SERVER_PORT='8250'
ENV SENZING_API_SERVER_SKIP_ENGINE_PRIMING='true'
ENV SENZING_API_SERVER_SKIP_STARTUP_PERF='true'
ENV SENZING_DATA_MART_SQLITE_DATABASE_FILE=/tmp/datamart
ENV SENZING_ENGINE_CONFIGURATION_JSON='{"PIPELINE": {"CONFIGPATH": "/etc/opt/senzing", "LICENSESTRINGBASE64": "", "RESOURCEPATH": "/opt/senzing/g2/resources", "SUPPORTPATH": "/opt/senzing/data"}, "SQL": {"CONNECTION": "sqlite3://na:na@/tmp/sqlite/G2C.db"}}'

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/demo-quickstart"]
