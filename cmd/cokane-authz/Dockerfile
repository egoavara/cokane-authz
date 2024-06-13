####################################################
# 빌드 단계, Go 어플리케이션을 빌드한다.
FROM golang:latest as builder

WORKDIR /usr/src/app

# 의존성 설치
COPY go.mod go.sum ./
RUN go mod download

# 소스코드 복사
COPY . .
RUN go mod tidy

# 빌드
RUN go build -v -o cokane-authz ./cmd/cokane-authz

####################################################
# 실행 단계, 빌드된 Go 어플리케이션을 실행한다.
FROM alpine:3.20 as runner

COPY --from=builder /usr/src/app/cokane-authz /usr/src/app/cokane-authz
ENTRYPOINT ["/usr/src/app/cokane-authz"]