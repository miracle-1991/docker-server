FROM golang:1.18

#install cron
RUN apt-get update && apt-get -y install cron && apt-get -y install vim

#set workdir
WORKDIR /app

#copy go module
COPY go.mod go.sum ./
RUN go mod download

#copy code
COPY . .
WORKDIR /app/cmd
RUN go build -o my-flight-service

#add cron task
ADD ./cmd/crontabfile /etc/cron.d/update-cron
RUN chmod 0644 /etc/cron.d/update-cron
RUN crontab /etc/cron.d/update-cron

# start
CMD ["sh", "-c", "cron && ./my-flight-service"]