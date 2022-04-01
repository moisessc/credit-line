FROM golang:alpine AS build

RUN apk add --update git
WORKDIR /go/credit-line-api
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/credit-line-api cmd/credit-line-api/main.go

FROM scratch
COPY --from=build /go/bin/credit-line-api /go/bin/credit-line-api
ENTRYPOINT ["/go/bin/credit-line-api"]