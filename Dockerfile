FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apt-get update -y && \
    apt-get install -y curl tar bzip2 libgl1-mesa-glx libglu1-mesa libgtk-3-dev libdbus-1-dev

## Install PrusaSlicer
RUN curl -L -o prusa-slicer.tar.bz2 https://github.com/prusa3d/PrusaSlicer/releases/download/version_2.6.1/PrusaSlicer-2.6.1+linux-x64-GTK3-202309060711.tar.bz2 && \
    mkdir -p /usr/local/bin/prusa-slicer && \
    tar -xjvf prusa-slicer.tar.bz2 -C /usr/local/bin/prusa-slicer --strip-components=1 && \
    chmod +x /usr/local/bin/prusa-slicer/bin/prusa-slicer && \
    rm -rf prusa-slicer.tar.bz2

ENV PATH="/usr/local/bin/prusa-slicer/bin:$PATH"
ENV LD_LIBRARY_PATH="/usr/lib/x86_64-linux-gnu:$LD_LIBRARY_PATH"

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/cloudslicer/

EXPOSE 8080

CMD ["./cloudslicer"]

