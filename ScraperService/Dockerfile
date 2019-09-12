FROM golang
MAINTAINER sourav241196@gmail.com
WORKDIR /app
ENV GO111MODULE=on
COPY ./go.mod ./go.sum ./
RUN go mod tidy
RUN go mod download
RUN go mod verify
COPY ./ .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./ScraperService/bin/scraper ./ScraperService
CMD ["./ScraperService/bin/scraper"]

FROM alpine:latest
COPY --from=0 /app/ScraperService/bin/scraper .
COPY --from=0 /app/dump ./dump
CMD ["./scraper"]