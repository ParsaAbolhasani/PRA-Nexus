# مرحله 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# نصب وابستگی‌های سیستمی
RUN apk add --no-cache git

# کپی فایل‌های go.mod و go.sum
COPY go.mod go.sum ./
RUN go mod download

# کپی کل سورس
COPY . .

# ساخت باینری
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pra-exchange .

# مرحله 2: اجرا (تصویر نهایی کوچک)
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# کپی باینری از مرحله قبل
COPY --from=builder /app/pra-exchange .
COPY --from=builder /app/.env .env

# کپی فایل Swagger (اختیاری)
COPY --from=builder /app/api/docs/swagger.yaml ./api/docs/

# اکسپوز پورت API
EXPOSE 8080

# اجرا
CMD ["./pra-exchange"]