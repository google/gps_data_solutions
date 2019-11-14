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

package split_file

import (
    "cloud.google.com/go/storage"
    "context"
    "io/ioutil"
    "os"
    "path"
    "time"
    "fmt"
    "log"
    "strconv"
    "strings"
)

type GCSEvent struct {
    Bucket         string    `json:"bucket"`
    Name           string    `json:"name"`
    Metageneration string    `json:"metageneration"`
    ResourceState  string    `json:"resourceState"`
    TimeCreated    time.Time `json:"timeCreated"`
    Updated        time.Time `json:"updated"`
}

func readLines(ctx context.Context, client *storage.Client, event GCSEvent) ([]string, error) {
    rc, err := client.Bucket(event.Bucket).Object(event.Name).NewReader(ctx)
    if err != nil {
        return nil, err
    }
    defer rc.Close()

    data, err := ioutil.ReadAll(rc)
    if err != nil {
        return nil, err
    }

    var lines []string
    lines = strings.Split(string(data), "\n")
    return lines, nil
}

func writeLines(ctx context.Context, client *storage.Client, lines []string, bucket string, file string) error {
    wc := client.Bucket(bucket).Object(file).NewWriter(ctx)
    if _, err := wc.Write([]byte(strings.Join(lines[:], "\n"))); err != nil {
        log.Printf("createFile: unable to write data to bucket %q, file %q: %v", bucket, file, err)
        return err
    }
    if err := wc.Close(); err != nil {
        log.Printf("createFile: unable to close bucket %q, file %q: %v", bucket, file, err)
        return err
    }
    return nil
}

func getEnv(key string, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func SplitFile(ctx context.Context, e GCSEvent) error {
    batchSize, _ := strconv.Atoi(getEnv("BATCH_SIZE", "1000"))
    outputBucket := getEnv("GCS_BUCKET_OUTPUT", "")
    outputDirectory := getEnv("GCS_DIRECTORY_OUTPUT", "")
    deleteFile, _ := strconv.ParseBool(getEnv("DELETE_FILE_AFTER_PROCESSING", "False"))

    p := path.Dir(e.Name)
    fn := strings.TrimPrefix(e.Name, fmt.Sprintf("%s/", p))
    fnc := strings.TrimSuffix(fn, path.Ext(fn))

    client, err := storage.NewClient(ctx)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    index := 0
    lines, err := readLines(ctx, client, e)
    if err != nil {
        log.Fatalf("Error reading file: %s", err)
    }
    for i := 0; i < len(lines); i += batchSize {
        j := i + batchSize
        if j > len(lines) {
            j = len(lines)
        }
        outputFile := path.Join(outputDirectory, fmt.Sprintf("%s-%d.%s", fnc, index, path.Ext(fn)))
        if err := writeLines(ctx, client, lines[i:j], outputBucket, outputFile); err != nil {
            log.Fatalf("Error writing file: %s", err)
            } else {
            log.Printf("File written with %d lines: %s", batchSize, path.Join(outputBucket, outputFile))
        }
        index++
    }

    log.Printf("Delete file: %t", deleteFile)
    if deleteFile {
        o := client.Bucket(e.Bucket).Object(e.Name)
        if err := o.Delete(ctx); err != nil {
            return err
        } else{
        log.Printf("File %s deleted.", path.Join(e.Bucket, e.Name))
        }
    }

    return nil
}
