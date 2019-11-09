# Start from the latest golang base image
FROM golang:1.12.0-alpine3.9


# We create an /app directory within our image that will hold our application source files
RUN mkdir /auction-system

# Add everything in root dir to the below folder
ADD . /auction-system

# cd into Work directory
WORKDIR /auction-system

RUN go get

RUN go build -o main .



# Command to run the executable
CMD ["./auction-system"]