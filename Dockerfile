FROM golang:1.13.1-stretch AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mpulse_exporter .

FROM golang:1.13.1-stretch AS app
USER 59000:59000
COPY --from=build /src/mpulse_exporter /mpulse_exporter
EXPOSE 8080
ENTRYPOINT [ "/mpulse_exporter" ]
