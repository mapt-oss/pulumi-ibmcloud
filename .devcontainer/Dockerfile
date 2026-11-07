FROM jetpackio/devbox:latest@sha256:7a9a41657bbfcb116d55811e4bfafd2c57bdaeb94557edb0308fe80b67831df2

# Installing your devbox project
WORKDIR /code
COPY devbox.json devbox.json
COPY devbox.lock devbox.lock
RUN sudo chown -R "${DEVBOX_USER}:${DEVBOX_USER}" /code


RUN devbox run -- echo "Installed Packages."

RUN devbox shellenv --init-hook >> ~/.profile
