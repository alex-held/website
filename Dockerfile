FROM alpine:latest as base
MAINTAINER Alexander Held 'contact@alexheld.io'
LABEL REPOSITORY="https://github.com/alex-held/website"

ENV PORT=8080 \
    GIN_MODE=release

EXPOSE $PORT

FROM golang:alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -a -o main .

FROM base as final

COPY --from=build /src/main /app/main
COPY --from=build /static/assets /app/assets
COPY --from=build /templates /app/templates

RUN ls /
RUN ls /app
RUN ls /app/templates
RUN ls /app/assets

CMD ["/app/main"]




# RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# RUN mkdir -p /api
# WORKDIR /api
# COPY --from=builder /api/app .
# COPY --from=builder /api/test.db .


#NTRYPOINT ["./app"]
