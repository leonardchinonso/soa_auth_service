FROM golang:alpine3.18 as builder

WORKDIR /soa_auth_service

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o soa_auth_service

RUN ls -la

FROM golang:alpine3.18

COPY --from=builder /soa_auth_service/soa_auth_service ./
COPY --from=builder /soa_auth_service/config/dev.env ./config/

EXPOSE 8080

CMD ["./soa_auth_service"]

