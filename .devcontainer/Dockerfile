FROM jetpackio/devbox:latest@sha256:3446246d7fbadd0fb0fa53f434593d01b5711148e7608accee74ac448137169e

# Installing your devbox project
WORKDIR /code
COPY devbox.json devbox.json
COPY devbox.lock devbox.lock
RUN sudo chown -R "${DEVBOX_USER}:${DEVBOX_USER}" /code


RUN devbox run -- echo "Installed Packages."

RUN devbox shellenv --init-hook >> ~/.profile
