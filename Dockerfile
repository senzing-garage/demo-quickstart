# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_GO_BUILDER=golang:1.21.4-bullseye
# FIXME: ARG IMAGE_FINAL=senzing/senzingapi-runtime:3.8.0
ARG IMAGE_FINAL=senzing/senzingapi-runtime:staging

# -----------------------------------------------------------------------------
# Stage: senzingapi_runtime
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} as senzingapi_runtime

# -----------------------------------------------------------------------------
# Stage: go_builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_GO_BUILDER} as go_builder
ENV REFRESHED_AT=2023-10-03
LABEL Name="senzing/demo-quickstart-builder" \
      Maintainer="support@senzing.com" \
      Version="0.2.0"

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
      Version="0.2.0"

# Copy files from prior stage.

COPY --from=go_builder "/output/linux-amd64/demo-quickstart" "/app/demo-quickstart"

# Runtime environment variables.

ENV LD_LIBRARY_PATH=/opt/senzing/g2/lib/

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/demo-quickstart"]
