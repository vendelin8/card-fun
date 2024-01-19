FROM golang:alpine
RUN apk add go-task
EXPOSE 3000
WORKDIR /app
COPY . .
CMD ["go-task", "run"]