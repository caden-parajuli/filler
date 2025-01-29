# Nix builder
FROM nixos/nix:latest AS builder

# Copy our source and setup our working dir.
COPY . /tmp/build
WORKDIR /tmp/build

RUN nix-shell -p busybox --command "addgroup -S nginx && adduser -S -s /bin/false nginx -G nginx"

EXPOSE 80

# Build our Nix environment
RUN nix \
    --extra-experimental-features "nix-command flakes" \
    # --option filter-syscalls false \
    build .#default

CMD ["/tmp/build/result/bin/default"]

