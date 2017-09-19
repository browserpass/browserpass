FROM golang:1.8-alpine3.6

ENV CGO_ENABLED=0 \
    APP_PATH="$GOPATH/src/github.com/dannyvankooten/browserpass"
# New ENV statement as this depends on APP_PATH
ENV PATH="$PATH:$APP_PATH/node_modules/.bin/"

RUN apk add --no-cache bash git make tar yarn zip && \
    go get -u github.com/golang/dep/cmd/dep && \
    mkdir -p $APP_PATH && \
    ln -s $APP_PATH /

WORKDIR $APP_PATH

ENTRYPOINT ["make"]
