// Copyright 2019 Google LLC
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     https://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bq_to_gcs

import (
    "cloud.google.com/go/bigquery"
    "context"
    "encoding/json"
    "fmt"
    "log"
    "errors"
    "net/http"
    "os"
)

type jobConfig struct {
    Project	string	`json:"project"`
    Dataset	string	`json:"dataset"`
    Table	string	`json:"table"`
    GcsURI	string	`json:"gcsURI"`
}

type PubSubMessage struct {
    Data []byte `json:"data"`
}

// exportTableAsJSONHandler starts the BigQuery job to export the table to Cloud Storage
func exportTableAsJSONHandler(ctx context.Context, client *bigquery.Client, conf jobConfig) error {
    gcsRef := bigquery.NewGCSReference(conf.GcsURI)
    gcsRef.DestinationFormat = bigquery.JSON
    extractor := client.DatasetInProject(conf.Project, conf.Dataset).Table(conf.Table).ExtractorTo(gcsRef)
    job, err := extractor.Run(ctx)
    if err != nil {
        return err
    }
    status, err := job.Wait(ctx)
    if err != nil {
        return err
    }
    if err := status.Err(); err != nil {
        return err
    }
    return nil
}

// Populates config var with environment variable data.
func setConfigFromEnvVars(b *jobConfig) error {
    if b.Project = os.Getenv("GCP_PROJECT"); b.Project == "" {
        return errors.New("GCP_PROJECT ENV NOT SET")
    }
    if b.Dataset = os.Getenv("BQ_DATASET"); b.Dataset == "" {
        return errors.New("BQ_DATASET ENV NOT SET")
    }
    if b.Table = os.Getenv("BQ_TABLE"); b.Table == "" {
        return errors.New("BQ_TABLE ENV NOT SET")
    }
    if b.GcsURI = os.Getenv("GCS_URI"); b.GcsURI == "" {
        return errors.New("GCS_URI ENV NOT SET")
    }
    return nil
}

// ExportTableAsJSONHttp is the entry point for the http trigger.
func ExportTableAsJSONHttp(w http.ResponseWriter, r *http.Request) {
    var b jobConfig
    if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
        log.Println("EMPTY CONFIG BODY, CHECKING ENVIROMENT VARIABLES: ")
        if err := setConfigFromEnvVars(&b); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }
    client, err := bigquery.NewClient(r.Context(), b.Project)
    if err != nil {
        log.Println("ERROR CREATING BQ CLIENT")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    if err := exportTableAsJSONHandler(r.Context(), client, b); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Printf("EXPORTED:\n%+v\n", b)
}

// ExportTableAsJSONPubSub is the entry point for the PubSub trigger.
func ExportTableAsJSONPubSub(ctx context.Context, m PubSubMessage) error {
    var b jobConfig
    if err := setConfigFromEnvVars(&b); err != nil {
        return fmt.Errorf("CONFIG ERR: %v", err)
    }
    client, err := bigquery.NewClient(ctx, b.Project)
    if err != nil {
        return fmt.Errorf("BQ ERROR: %v", err)
    }
    if err := exportTableAsJSONHandler(ctx, client, b); err != nil {
        return fmt.Errorf("EXPORT ERROR: %v", err)
    }
    fmt.Printf("EXPORTED:\n%+v\n", b)
    return nil
}