# Stage 1: Build the Go application
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cloudslicer ./cmd/cloudslicer/

# Stage 2: Create an image for PrusaSlicer and its dependencies
FROM debian:bullseye-slim AS prusa-slicer

RUN apt-get update -y && \
    apt-get install -y curl tar bzip2 libgl1-mesa-glx libglu1-mesa libgtk-3-dev libdbus-1-dev && \
    mkdir -p /app/temp/gcode && \
    mkdir -p /app/temp/model && \
    mkdir -p /app/temp/config && \
    curl -L -o prusa-slicer.tar.bz2 https://github.com/prusa3d/PrusaSlicer/releases/download/version_2.6.1/PrusaSlicer-2.6.1+linux-x64-GTK3-202309060711.tar.bz2 && \
    mkdir -p /usr/local/bin/prusa-slicer && \
    tar -xjvf prusa-slicer.tar.bz2 -C /usr/local/bin/prusa-slicer --strip-components=1 && \
    chmod +x /usr/local/bin/prusa-slicer/bin/prusa-slicer && \
    rm -rf prusa-slicer.tar.bz2


# Stage 3: Create the final image
FROM debian:bullseye-slim

WORKDIR /app

# Copy the built Go binary from Stage 1 && the presets directory
COPY --from=builder /app/cloudslicer .

# Copy PrusaSlicer and its dependencies from Stage 2
COPY --from=prusa-slicer / /

ENV PATH="/usr/local/bin/prusa-slicer/bin:$PATH"
ENV LD_LIBRARY_PATH="/usr/lib/x86_64-linux-gnu:$LD_LIBRARY_PATH"

# Expose the required port
EXPOSE 8080

# Command to run the application
CMD ["./cloudslicer"]

