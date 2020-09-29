FROM golang:1.15 as builder

ENV APP_NAME IPLocator
ENV APP_USER app
ENV APP_HOME /go/src/github.com/Tomer-Alony/IPLocator

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME

WORKDIR $APP_HOME
USER $APP_USER
COPY src/ .

RUN go mod download
RUN go mod verify
RUN go build -o $APP_NAME

FROM debian:buster

ENV APP_USER app
ENV BIN_HOME /go/src/github.com/Tomer-Alony/IPLocator
ENV APP_HOME /go/src/$APP_NAME
ENV PORT 8080

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

COPY --chown=0:0 --from=builder $BIN_HOME/$APP_NAME $APP_HOME
RUN ls -lh ./

EXPOSE $PORT
USER $APP_USER
CMD ["./IPLocator"]