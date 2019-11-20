# Collection of Cloud Functions

**Please note: this is not an officially supported Google product.**

Function          | Description  | Language
----------------- | ------------ | ------------------
[Vision API](vision_api/) | Send files to the Cloud Vision API | NodeJS
[Split File](split_file/) | Split GCS stored file into smaller chuncks | Go
[BQ to GCS](bg_to_gcs/) | Export BQ table to Cloud Storage | Go

# Vision API
This cloud function triggers on new objects stored on Google Cloud Storage. Each file will be send to the Cloud Vision API to extract some data. The Vision API results are stored in BigQuery. The following settings can be configured by setting environment variables:

Variable | Description | Default
--------| -----------| ---------
BQ_DATASET | Output BigQuery Dataset | 
BQ_TABLE | Output BigQuery Table | 

# Split File
Splits file stored to Google Cloud Storage into smaller files. The cloud function triggers on a new object created in a bucket. The following settings can be configured by setting environment variables:

Variable | Description | Default
--------| -----------| ---------
BATCH_SIZE | Number of lines in the output file. | 1000
GCS_BUCKET_OUTPUT | Output cloud storage bucket | 
GCS_DIRECTORY_OUTPUT | Output directory inside the output bucket. | 
DELETE_FILE_AFTER_PROCESSING | If set to True it removes the input file after processing. | False

# BQ to GCS
Exports a BigQuery table to Google Cloud Storage, can be configured to trigger on http request or PubSub. The following settings can be configured by setting environment variables:

Variable | Description | Default
--------| -----------| ---------
BQ_DATASET | BigQuery Dataset | 
BQ_TABLE | BigQuery Table to export | 
GCS_URI | Google Cloud Storage path. e.g. gs//bucket_name/file_name.json | 