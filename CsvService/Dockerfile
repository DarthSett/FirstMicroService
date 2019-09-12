FROM golang
MAINTAINER sourav241196@gmail.com
WORKDIR /app
ENV GO111MODULE=on
COPY ./go.mod ./go.sum ./
RUN go mod tidy
RUN go mod download
RUN go mod verify
COPY ./ .
RUN CGO_ENABLED=0 go build -v -o ./CsvService/bin/uploader ./CsvService
EXPOSE 4000
CMD ["./CsvService/bin/uploader"]

FROM alpine:latest
COPY --from=0 /app/CsvService/bin/uploader .
CMD ["./uploader"]