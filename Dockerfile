FROM golang:1.16.5-buster AS builder

WORKDIR /go/src/app


COPY . .

ENV CTX_TIMEOUT=100
ENV PORT=3030
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

RUN mv /etc/localtime localtime.backup
RUN ln -s /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

RUN #cp config/.env.example config/.env
RUN rm -rf vendor && \
        go mod tidy && \
        go mod vendor && \
        go build -o backend-service

FROM ubuntu:20.04


# Import the binary of sql-migrate, user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable.
COPY --from=builder /go/src/app /app
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Jakarta
ENV CTX_TIMEOUT=100
ENV PORT=3030

# Use an unprivileged user.
USER appuser:appuser

EXPOSE 3030
# Run the hello binary.
ENTRYPOINT ["/app/backend-service"]
