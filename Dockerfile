ARG project=github.com/jannis-a/go-durak
ARG workdir=/go/src/${project}

FROM golang:1.12-alpine3.9 AS builder
ARG workdir
WORKDIR ${workdir}
COPY . .
RUN apk --no-cache add build-base git
RUN make build-prod

FROM scratch
ARG workdir
COPY --from=builder ${workdir}/_output/app .
CMD ["./app"]
