FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
COPY ./etc /app/etc
RUN go build -ldflags="-s -w" -o /app/website ./cmd/website/website.go
RUN go build -ldflags="-s -w" -o /app/deposit ./cmd/deposit/deposit.go
RUN go build -ldflags="-s -w" -o /app/inscribe ./cmd/inscribe/inscribe.go


# Targets

FROM scratch as website

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/website /app/website
COPY --from=builder /app/etc /app/etc

CMD ["./website", "-f", "etc/website-api.yaml"]



FROM scratch  as deposit

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/deposit /app/deposit
COPY --from=builder /app/etc /app/etc

CMD ["./deposit", "-f", "etc/website-api.yaml"]


FROM scratch  as inscribe

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/inscribe /app/inscribe
COPY --from=builder /app/etc /app/etc

CMD ["./inscribe", "-f", "etc/website-api.yaml"]
