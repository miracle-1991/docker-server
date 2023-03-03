FROM golang:latest

#install cron
RUN apt-get update && apt-get -y install cron

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
RUN echo "0 */6 * * * /app/my-flight-service > /dev/null 2>&1" >> /etc/crontab

# start
CMD ["sh", "-c", "cron && ./my-flight-service"]