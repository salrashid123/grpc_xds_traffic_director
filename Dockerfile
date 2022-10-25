FROM golang:1.19 as build

RUN apt-get update -y && apt-get install -y build-essential wget unzip curl git


RUN curl -OL https://github.com/google/protobuf/releases/download/v3.19.0/protoc-3.19.0-linux-x86_64.zip && \
    unzip protoc-3.19.0-linux-x86_64.zip -d protoc3 && \
    mv protoc3/bin/* /usr/local/bin/ && \
    mv protoc3/include/* /usr/local/include/


ENV GO111MODULE=on
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /app

ADD . /app

RUN go mod download


RUN /usr/local/bin/protoc --go_out=.  --go_opt=paths=source_relative  --descriptor_set_out=echo/echo.proto.pb     --go-grpc_out=. --go-grpc_opt=paths=source_relative   echo/echo.proto


RUN export GOBIN=/app/bin && go install server/grpc_server.go
RUN export GOBIN=/app/bin && go install client/grpc_client.go


FROM gcr.io/distroless/base
COPY --from=build /app/certs /certs/
COPY --from=build /app/bin/grpc_server /
COPY --from=build /app/bin/grpc_client /

EXPOSE 50051

ENTRYPOINT ["/grpc_server", "--grpcport", ":50051", "--hcaddress", ":50050", "-tlsCert", "certs/grpc.crt", "-tlsKey", "certs/grpc.key"]
#ENTRYPOINT ["/grpc_client", "--host",  "grpc.domain.com:50051", "-tlsCert", "certs/tls-ca-chain.pem"]
