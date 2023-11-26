FROM golang

ADD . /usr/local/src
WORKDIR /usr/local/src

RUN go build -o /server server/server.go
RUN go build -o /proxy proxy/proxy.go

CMD ["/server --port 8080 --dbpassword my_postgre --runviadocker true"]
CMD ["/proxy --port 3000 --runviadocker true"]