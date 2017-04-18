# presigner [ ![Codeship Status for minodisk/presigner](https://img.shields.io/codeship/7e793c10-01bb-0135-39d5-52a787a130cb/master.svg?style=flat)](https://app.codeship.com/projects/212925) [![Go Report Card](https://goreportcard.com/badge/github.com/minodisk/presigner)](https://goreportcard.com/report/github.com/minodisk/presigner) [![codecov](https://codecov.io/gh/minodisk/presigner/branch/master/graph/badge.svg)](https://codecov.io/gh/minodisk/presigner) [![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat)](https://godoc.org/github.com/minodisk/presigner) [![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

Pre-signed URL publisher to upload files directly to Google Cloud Storage.

## Usage

### 1. Setup bucket

#### [Required] Create a bucket:

```sh
gsutil mb gs://example-bucket
```

#### [Optional] Set CORS to the bucket:

To upload from browser with XHR. See [Cross-Origin Resource Sharing (CORS) - Cloud Storage — Google Cloud Platform](https://cloud.google.com/storage/docs/cross-origin).

```sh
gsutil cors set example-cors.json gs://example-bucket
```

#### [Optional] Set default object ACL to the bucket:

To make objects accessible from any users. See [defacl - Get, set, or change default ACL on buckets - Cloud Storage — Google Cloud Platform](https://cloud.google.com/storage/docs/gsutil/commands/defacl#ch).

```sh
gsutil defacl ch -u AllUsers:R gs://example-bucket
```

### 2. Generate private key

[Generate JSON private key for service account](https://cloud.google.com/storage/docs/authentication#generating-a-private-key).

### 3. Run

```sh
presigner -account /path/to/private-key.json -bucket bucket-a -port 80
```

#### Options:

```sh
presigner -help
```

## HTTP(S) API

### Publish signed-URL

#### Request Body:

```json
{
  "bucket": "example",
  "content_type": "image/jpeg",
  "md5": "XXXXXXXXXXXXXXXXXXXX"
}
```

- `bucket`: Bucket name to upload file.
- `content_type`: Content type of the file.
- `md5`: [Optional] MD5 checksum of the file.

#### Response Body:

```json
{
  "signed_url": "https://storage.googleapis.com/presigner/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX?Expires=1492484353&GoogleAccessId=example%40xxx.iam.gserviceaccount.com&Signature=...",
  "file_url": "https://example.storage.googleapis.com/XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
}
```

- `signed_url`: Endpoint of uploading file.
- `file_url`: URL of uploaded file.

##### Error case

```json
{
  "error": ""
}
```

- `error`: Reason of error.

### Upload file

`PUT` to `signed_url`, and write contents of the file to request body.

#### HTTP Header:

```http
Content-Type: image/jpeg
Content-Disposition: attachement; filename="example.jpg"
```

- `Content-Type`: Content type of the file. Same as `content_type` specified when publishing signed-URL.
- `Content-Disposition`: [Optional] Name given when the file is downloaded.
