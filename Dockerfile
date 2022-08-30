FROM golang:latest

#设置工作目录
WORKDIR $GOPATH/src/chat/src/main

#将服务器的go工程代码加入到docker容器中
ADD . $GOPATH/src/chat

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV GOPATH=$HOME/go
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/go/bin

#go构建可执行文件
RUN go build -o ../../bin/chat

EXPOSE 80

#最终运行docker命令
ENTRYPOINT ["../../bin/chat"]
