# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_BUILDER=golang:1.22.4-bullseye
ARG IMAGE_FINAL=senzing/senzingapi-runtime-staging:latest

# -----------------------------------------------------------------------------
# Stage: senzingapi_runtime
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} AS senzingapi_runtime

# -----------------------------------------------------------------------------
# Stage: builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_BUILDER} AS builder
ENV REFRESHED_AT=2024-07-01
LABEL Name="senzing/go-builder" \
      Maintainer="support@senzing.com" \
      Version="0.1.0"

# Run as "root" for system installation.

USER root

# Install packages via apt-get.

RUN apt-get update \
 && apt-get -y install \
      python3 \
      python3-dev \
      python3-pip \
      python3-venv \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

# Create and activate virtual environment.

RUN python3 -m venv /app/venv
ENV PATH="/app/venv/bin:$PATH"

# Install packages via PIP.

COPY requirements.txt .
RUN pip3 install --upgrade pip \
 && pip3 install -r requirements.txt \
 && rm requirements.txt \
 && pip3 uninstall -y setuptools pip

# Copy local files from the Git repository.

COPY ./rootfs /
COPY . ${GOPATH}/src/demo-quickstart

# Copy files from prior stage.

COPY --from=senzingapi_runtime  "/opt/senzing/er/lib/"   "/opt/senzing/er/lib/"
COPY --from=senzingapi_runtime  "/opt/senzing/er/sdk/c/" "/opt/senzing/er/sdk/c/"

# Set path to Senzing libs.

ENV LD_LIBRARY_PATH=/opt/senzing/er/lib/

# Build go program.

WORKDIR ${GOPATH}/src/demo-quickstart
RUN make build

# Copy binaries to /output.

RUN mkdir -p /output \
 && cp -R ${GOPATH}/src/demo-quickstart/target/*  /output/

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} AS final
ENV REFRESHED_AT=2024-07-01
LABEL Name="senzing/demo-quickstart" \
      Maintainer="support@senzing.com" \
      Version="0.0.1"
HEALTHCHECK CMD ["/app/healthcheck.sh"]
USER root

# Install packages via apt-get.

RUN export STAT_TMP=$(stat --format=%a /tmp) \
 && chmod 777 /tmp \
 && apt-get update \
 && apt-get -y install \
      gnupg2 \
      jq \
      libodbc1 \
      postgresql-client \
      supervisor \
      unixodbc \
 && chmod ${STAT_TMP} /tmp \
 && rm -rf /var/lib/apt/lists/*

# Install Java-11.

RUN mkdir -p /etc/apt/keyrings \
 && wget -O - https://packages.adoptium.net/artifactory/api/gpg/key/public > /etc/apt/keyrings/adoptium.asc

RUN echo "deb [signed-by=/etc/apt/keyrings/adoptium.asc] https://packages.adoptium.net/artifactory/deb $(awk -F= '/^VERSION_CODENAME/{print$2}' /etc/os-release) main" >> /etc/apt/sources.list

RUN export STAT_TMP=$(stat --format=%a /tmp) \
 && chmod 777 /tmp \
 && apt-get update \
 && apt-get -y install \
      temurin-11-jdk \
      python3-venv \
      curl \
 && chmod ${STAT_TMP} /tmp \
 && rm -rf /var/lib/apt/lists/*

# Copy files from repository.

COPY ./rootfs /

# Copy files from prior stage.

COPY --from=builder /output/linux/demo-quickstart /app/demo-quickstart
COPY --from=builder /app/venv /app/venv

# Prepare jupyter lab environment

RUN mkdir -p /.local/share /notebooks

# Run as non-root container

USER 1001

# Activate virtual environment.

ENV VIRTUAL_ENV=/app/venv
ENV PATH="/app/venv/bin:${PATH}"

# Runtime environment variables.

ENV LD_LIBRARY_PATH=/opt/senzing/er/lib/
ENV SENZING_API_SERVER_ALLOWED_ORIGINS='*'
ENV SENZING_API_SERVER_BIND_ADDR='all'
ENV SENZING_API_SERVER_ENABLE_ADMIN='true'
ENV SENZING_API_SERVER_PORT='8250'
ENV SENZING_API_SERVER_SKIP_ENGINE_PRIMING='true'
ENV SENZING_API_SERVER_SKIP_STARTUP_PERF='true'
ENV SENZING_DATA_MART_SQLITE_DATABASE_FILE=/tmp/datamart
ENV SENZING_ENGINE_CONFIGURATION_JSON='{"PIPELINE": {"CONFIGPATH": "/etc/opt/senzing", "LICENSESTRINGBASE64": "", "RESOURCEPATH": "/opt/senzing/er/resources", "SUPPORTPATH": "/opt/senzing/data"}, "SQL": {"CONNECTION": "sqlite3://na:na@/tmp/sqlite/G2C.db"}}'

# Runtime execution.

WORKDIR /app
CMD ["/usr/bin/supervisord", "--nodaemon"]
