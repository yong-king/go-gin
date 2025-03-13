## 使用最新的 Golang 镜像
#FROM golang:1.23.0
#
## 设置工作目录为 /app
#WORKDIR /app
#
## 复制项目文件到 /app
#COPY go.mod go.sum ./
#COPY . /app
#
## 下载依赖
#RUN go mod tidydocker ps -as
#
## **确保 `gin-blog` 可执行文件生成到 `/app` 目录**
#RUN go build -o /app/gin-blog ./gin-blog
#
## 暴露端口
#EXPOSE 8000
#
## **确保执行正确的可执行文件**
#ENTRYPOINT ["/app/gin-blog/gin-blog"]


# 使用 scratch 镜像
# 选择最小化的 scratch 镜像
FROM scratch

# 设置工作目录
WORKDIR /app

# 复制编译好的可执行文件
COPY gin-blog/gin-blog /app/gin-blog

# 复制配置文件
COPY gin-blog/conf /app/conf

# 运行程序
ENTRYPOINT ["/app/gin-blog"]



