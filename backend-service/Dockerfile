# Utilize Docker's Multi Staged building to minimalize image size
FROM golang:alpine AS builder

# Workdir to make it more easy to transfer file
WORKDIR /build
# Transfer the package file
COPY hon-producer/go.mod hon-producer/go.sum ./

# Transfer the entire source code of jhon-producer
COPY hon-producer/ .
# Transfer the shared folder
COPY shared ./shared

# Make the /shared can be accessed via replacing dependency
RUN go mod edit -replace=github.com/jirbthagoras/hon/shared=/build/shared

# Download the dependencies
RUN go mod download
# Build the output binary file
RUN go build -o /out .

# The last stage of multi staged building
FROM alpine

WORKDIR /out

# Copy the executable from builder
COPY --from=builder /out /out

EXPOSE 3000

# Run the executable
CMD ["/build"]
