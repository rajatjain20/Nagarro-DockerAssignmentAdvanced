# Base golang image with latest go version
FROM golang:1.23.1

# set working directory inside container
WORKDIR /go/app/

# copy everything from current directory from host to workimng directory inside container
# this will ignore few files specified in .dockerignore file
COPY ./src .

# set the TimeZone env variable, so that the logs will have IST time
ENV TZ="Asia/Kolkata"

# build the go code inside container
# this will generate executale "datainputapp"
RUN go build -v -o ./ .

# Expose the port
EXPOSE 3400

# command to execute the application
CMD [ "./datainputapp" ]
