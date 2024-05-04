FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o score-keep-api .

FROM alpine:3.14
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=build /app/score-keep-api .

EXPOSE 4000
CMD ["/app/score-keep-api"]