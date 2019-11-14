#!/bin/bash

# Copyright 2019 Google LLC.
# SPDX-License-Identifier: Apache-2.0

gcloud functions deploy ExportTableAsJSONPubSub \
    --runtime go111 \
    --trigger-topic "TOPIC" \
    --set-env-vars BQ_DATASET="DATASET",BQ_TABLE="TABLE",GCS_URI="GCS_URI"