FROM golang:1.24.4 AS build
COPY . ./
RUN go install github.com/a-h/templ/cmd/templ@v0.3.857
RUN templ generate
RUN CGO_ENABLED=1 go build -o "/bin/scnorionplus-console" .

FROM debian:latest
COPY --from=build /bin/scnorionplus-console /bin/scnorionplus-console
COPY ./assets /bin/assets
RUN apt-get update
RUN apt install -y ca-certificates
EXPOSE 1323
EXPOSE 1324
WORKDIR /bin
ENTRYPOINT ["/bin/scnorionplus-console"]