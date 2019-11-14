/**
 * Copyright 2019 Google LLC
 * 
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * 
 *     https://www.apache.org/licenses/LICENSE-2.0
 * 
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
const { GoogleAuth } = require('google-auth-library');
const { BigQuery } = require('@google-cloud/bigquery');

const FEATURES = [
    'FACE_DETECTION',
    'LANDMARK_DETECTION',
    'LOGO_DETECTION',
    'LABEL_DETECTION',
    'TEXT_DETECTION',
    'SAFE_SEARCH_DETECTION',
    'IMAGE_PROPERTIES',
]

const MAX_RESULTS = process.env.MAX_VISION_RESULTS;

/**
 * Sends the stored file on GCS to the Vision API and returns the first response.
 * @param {object} file The event payload from the Cloud Function entry point.
 */
async function getImageInfo(file) {
    const auth = new GoogleAuth({
        scopes: [
            'https://www.googleapis.com/auth/cloud-platform',
            'https://www.googleapis.com/auth/cloud-vision'
        ]
    });

    const vision = await auth.getClient();
    const url = `https://vision.googleapis.com/v1p4beta1/images:annotate`;
    const request = {
        requests: [
            {
                image: {
                    source: {
                        gcsImageUri: `gs://${file.bucket}/${file.name}`
                    }
                },
                features: FEATURES.map((feature) => {
                    return { type: feature, maxResults: MAX_RESULTS };
                })
            }
        ]
    }

    const response = await vision.request({
        url,
        data: JSON.stringify(request),
        method: 'POST'
    });

    if (response.status === 200) {
        const row = response.data.responses[0];
        row.gcsPath = `${file.bucket}/${file.name}`;
        row.gcsBucket = file.bucket;
        row.gcsFile = file.name;
        row.gcsCreated = file.timeCreated;
        row.gcsUpdated = file.updated;
        return row;
    }
    return null;
}

/**
 * getVisionAPIResults is the Cloud function entry point.
 * Triggers on a Google Cloud Storage finalize event.
 * @param {object} data The event payload.
 * @param {object} context The event metadata.
 */
exports.getVisionAPIResults = async (data, context) => {
    const file = data;
    console.log(`[${context.eventId}] - Bucket: ${file.bucket}`);
    console.log(`[${context.eventId}] - File: ${file.name}`);
    console.log(`[${context.eventId}] - Created: ${file.timeCreated}`);
    console.log(`[${context.eventId}] - Updated: ${file.updated}`);
    const row = await getImageInfo(file);
    const bigquery = new BigQuery();
    await bigquery
        .dataset(process.env.BQ_DATASET)
        .table(process.env.BQ_TABLE)
        .insert([row], {
            ignore_unknown_values: true
        });
};