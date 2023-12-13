# syntax=docker/dockerfile:1

FROM golang:1.21 AS build-stage

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /gugcp

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/accesskey.json /accesskey.json
COPY --from=build-stage /gugcp /gugcp

ENV DB_HOST="34.101.147.33"
ENV DB_PORT="3306"
ENV DB_USERNAME="one_user"
ENV DB_PASSWORD="Rahasia88oke#"
ENV DB_NAME="one_db"
ENV MODE="production"

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["./gugcp"]