FROM golang:1.23 AS build

WORKDIR /knockbox/src
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /knockbox/bin/main

FROM gcr.io/distroless/static-debian12

COPY --from=build /knockbox/bin /

EXPOSE 9090

CMD ["/main"]