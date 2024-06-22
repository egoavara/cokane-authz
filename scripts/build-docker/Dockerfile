####################################################
# 빌드 단계, Go 어플리케이션을 빌드한다.
FROM golang:1.22.4-alpine3.20 as builder

WORKDIR /usr/src/cokane-authz

# 의존성 설치
COPY go.mod go.sum ./
RUN go mod download

# 소스코드 복사
COPY . .

# 빌드
RUN go mod tidy 
RUN go build -o cokane-authz ./cmd/cokane-authz

####################################################
# 실행 단계, 빌드된 Go 어플리케이션을 실행한다.
FROM alpine:3.20 as runner

COPY --from=builder /usr/src/cokane-authz/cokane-authz /usr/src/cokane-authz/cokane-authz

ENTRYPOINT ["/bin/sh" , "-c", "sleep Infinity"]