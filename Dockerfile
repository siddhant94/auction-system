# Start from the latest golang base image
FROM golang:1.12.0-alpine3.9


# We create an /app directory within our image that will hold our application source files
RUN mkdir /go/src/auction-system

# Add everything in root dir to the below folder
ADD . /go/src/auction-system

# cd into Work directory
WORKDIR /go/src/auction-system

RUN go get
#github.com/siddhant94/auction-system

RUN go build -o main .

RUN ls

# Command to run the executable
CMD ["./main"]