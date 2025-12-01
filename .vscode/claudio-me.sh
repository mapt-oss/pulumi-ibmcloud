#!/bin/bash

podman run -it --rm --user 0 --workdir /home/default/workdir \
        -v $PWD:/home/default/workdir:z \
        -v claudio-gcp:/home/default/.config/gcloud:Z \
        -e ANTHROPIC_VERTEX_PROJECT_ID=itpc-gcp-ai-eng-claude \
        -e ANTHROPIC_VERTEX_PROJECT_QUOTA=cloudability-it-gemini \
        quay.io/aipcc-cicd/claudio:v0.1.1

