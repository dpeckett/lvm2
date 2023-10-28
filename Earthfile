VERSION 0.7
FROM golang:1.21-bookworm
WORKDIR /workspace

tidy:
  LOCALLY
  RUN go mod tidy
  RUN go fmt ./...

lint:
  FROM golangci/golangci-lint:v1.54.2
  WORKDIR /workspace
  COPY . .
  RUN golangci-lint run --timeout 5m ./...

test:
  RUN apt update
  RUN apt install -y --no-install-recommends kmod lvm2 qemu-utils udev
  COPY +modules/modules /lib/modules
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  ENV SYSTEMD_IGNORE_CHROOT=1
  COPY ./testdata/udev/rules.d /etc/udev/rules.d
  RUN --privileged \
    mount -t devtmpfs none /dev \
    && mount -t devpts none /dev/pts \
    && /lib/systemd/systemd-udevd --daemon \
    && go test -coverprofile=coverage.out -v ./...
  SAVE ARTIFACT ./coverage.out AS LOCAL coverage.out

modules:
  LOCALLY
  SAVE ARTIFACT /lib/modules