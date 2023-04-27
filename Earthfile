VERSION 0.7
PROJECT tochemey/goakt


FROM tochemey/docker-go:1.20.1-0.7.0

# install the various tools to generate connect-go
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

# run a PR branch is created
#pr-pipeline:
#  PIPELINE
#  TRIGGER pr main
#  BUILD +lint
#  BUILD +local-test

# run on when a push to main is made
#main-pipeline:
#  PIPELINE
#  TRIGGER push main
#  BUILD +lint
#  BUILD +local-test

pbs:
    BUILD +internal-pb
    BUILD +protogen
    BUILD +sample-pb

test:
  BUILD +lint
  BUILD +local-test

code:

    WORKDIR /app

    # download deps
    COPY go.mod go.sum ./
    RUN go mod download -x

    # copy in code
    COPY --dir . ./

vendor:
    FROM +code

    RUN go mod vendor
    SAVE ARTIFACT /app /files

lint:
    FROM +vendor

    COPY .golangci.yml ./

    # Runs golangci-lint with settings:
    RUN golangci-lint run --timeout 10m


local-test:
    FROM +vendor

    WITH DOCKER --pull postgres:11
        RUN go test -mod=vendor ./... -race -v -coverprofile=coverage.out -covermode=atomic -coverpkg=./...
    END

    SAVE ARTIFACT coverage.out AS LOCAL coverage.out

internal-pb:
    # copy the proto files to generate
    COPY --dir protos/ ./
    COPY buf.work.yaml buf.gen.yaml ./

    # generate the pbs
    RUN buf generate \
            --template buf.gen.yaml \
            --path protos/internal/goakt

    # save artifact to
    SAVE ARTIFACT gen/goakt AS LOCAL internal/goakt

protogen:
    # copy the proto files to generate
    COPY --dir protos/ ./
    COPY buf.work.yaml buf.gen.yaml ./

    # generate the pbs
    RUN buf generate \
            --template buf.gen.yaml \
            --path protos/public/messages

    # save artifact to
    SAVE ARTIFACT gen/messages AS LOCAL messages

testprotogen:
    # copy the proto files to generate
    COPY --dir protos/ ./
    COPY buf.work.yaml buf.gen.yaml ./

    # generate the pbs
    RUN buf generate \
            --template buf.gen.yaml \
            --path protos/test/pb

    # save artifact to
    SAVE ARTIFACT gen gen AS LOCAL test/data

sample-pb:
    # copy the proto files to generate
    COPY --dir protos/ ./
    COPY buf.work.yaml buf.gen.yaml ./

    # generate the pbs
    RUN buf generate \
            --template buf.gen.yaml \
            --path protos/sample/pb

    # save artifact to
    SAVE ARTIFACT gen gen AS LOCAL examples/protos

compile-actor-cluster:
    COPY +vendor/files ./

    RUN go build -mod=vendor  -o bin/accounts ./examples/actor-cluster/k8s
    SAVE ARTIFACT bin/accounts /accounts

actor-cluster-image:
    FROM alpine:3.16.2

    WORKDIR /app
    COPY +compile-actor-cluster/accounts ./accounts
    RUN chmod +x ./accounts

    EXPOSE 50051
    EXPOSE 9000

    ENTRYPOINT ["./accounts"]
    SAVE IMAGE accounts:dev
