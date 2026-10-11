FROM jetpackio/devbox:latest@sha256:28f736f8a3568c47831038f66b7856dd375488cb099999299cb791d95aa9e82e

# Installing your devbox project
WORKDIR /code
COPY devbox.json devbox.json
COPY devbox.lock devbox.lock
RUN sudo chown -R "${DEVBOX_USER}:${DEVBOX_USER}" /code


RUN devbox run -- echo "Installed Packages."

RUN devbox shellenv --init-hook >> ~/.profile
