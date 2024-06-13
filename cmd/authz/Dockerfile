####################################################
# 빌드 단계, Go 어플리케이션을 빌드한다.
FROM golang:latest as builder

# 의존성 설치
COPY go.mod go.sum ./
RUN go mod download

# 소스코드 복사
COPY . .

# 빌드
RUN go build -o cluster ./cmd/cluster
####################################################
# 실행 단계, 빌드된 Go 어플리케이션을 실행한다.
FROM alpine:3.20 as runner

COPY --from=builder /app/cluster /app/cluster
ENTRYPOINT ["/app/cluster"]