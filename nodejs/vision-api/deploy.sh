#!/bin/bash

# Copyright 2019 Google LLC.
# SPDX-License-Identifier: Apache-2.0

gcloud functions deploy getVisionAPIResults \
    --runtime nodejs10 \
    --trigger-resource "CLOUD BUCKET" \
    --trigger-event google.storage.object.finalize