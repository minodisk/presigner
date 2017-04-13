# presigner [ ![Codeship Status for minodisk/presigner](https://app.codeship.com/projects/7e793c10-01bb-0135-39d5-52a787a130cb/status?branch=master)](https://app.codeship.com/projects/212925)

Pre-signed URL publisher to upload files directly to Google Cloud Storage.

## Usage

### Setup

#### Create a bucket:

Example:

```sh
gsutil create gs://example-bucket
```

#### Set CORS to the bucket:

To upload from browser with XHR. See [Cross-Origin Resource Sharing (CORS) - Cloud Storage — Google Cloud Platform](https://cloud.google.com/storage/docs/cross-origin).

Example:

```sh
gsutil cors set example-cors.json gs://example-bucket
```

#### Set default object ACL to the bucket:

To make objects accessible from any users. See [defacl - Get, set, or change default ACL on buckets - Cloud Storage — Google Cloud Platform](https://cloud.google.com/storage/docs/gsutil/commands/defacl#ch).

Example:

```sh
gsutil defacl ch -u AllUsers:R gs://example-bucket
```

### Run

1. Save Google Private Key as file.
2. Run presigner.

```sh
presigner -id example@project.iam.gserviceaccount.com -key /path/to/service-account.pem -bucket bucket-a -bucket bucket-b
```

#### More information:

```sh
presigner -help
```

## Upload file to GCS with pre-signed URL

TODO: Add link to GCS document
