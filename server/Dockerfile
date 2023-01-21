# Stage: Build
FROM golang:1.19-alpine3.17 AS build-env
WORKDIR /build
ADD . .
RUN go mod download
RUN go build -o /server .

# Stage: Run
FROM alpine:3.17
WORKDIR /goteam
COPY --from=build-env /server .
COPY --from=build-env /build/.env .
EXPOSE 8080
CMD ["/goteam/server"]
