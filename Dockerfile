FROM docker.io/golang:1.25 as build

RUN go env -w CGO_ENABLED=0

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ .
RUN go build -o /out/app ./cmd/app

FROM scratch

COPY migrations/ /migrations
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# This is for clarity when using top or htop
COPY --from=build /out/app /bin/app

EXPOSE 8080

CMD ["/bin/app"]