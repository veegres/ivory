FROM docker:29-dind

# openssh is Ivory's only transport to a docker platform host. The rest is what
# GET /node/platform/system/info shells out to: busybox answers a subset of
# those flags, and the call reports exit code 1 without musl-locales in
# particular, since `locale` is the last statement of its script and so owns
# the script's exit status.
RUN apk add --no-cache \
    openssh-server openssh-client \
    coreutils procps util-linux musl-locales pciutils net-tools

COPY start.sh /usr/local/bin/start.sh
RUN chmod +x /usr/local/bin/start.sh

ENTRYPOINT ["/usr/local/bin/start.sh"]
