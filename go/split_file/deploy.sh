#!/bin/bash

# Copyright 2019 Google LLC.
# SPDX-License-Identifier: Apache-2.0

gcloud functions deploy SplitFile \
    --runtime go111 \
    --trigger-resource "TRIGGER_GCS_BUCKET" \
    --trigger-event google.storage.object.finalize \
    --set-env-vars BATCH_SIZE="1000",GCS_BUCKET_OUTPUT="OUTPUT_BUCKET",GCS_DIRECTORY_OUTPUT="OUTPUT_DIRECTORY",DELETE_FILE_AFTER_PROCESSING=True
