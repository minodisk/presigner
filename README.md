# presigner [![Circle CI](https://circleci.com/gh/go-microservices/presigner/tree/master.svg?style=svg)](https://circleci.com/gh/go-microservices/presigner/tree/master)

Publisher of pre-signed URLs to upload files directly to Google Cloud Storage.

## Setup

1. Setup CORS to your bucket: See [Cross-Origin Resource Sharing (CORS) - Cloud Storage — Google Cloud Platform](https://cloud.google.com/storage/docs/cross-origin)
2. Setup default object ACL: See [defacl - Get, set, or change default ACL on buckets - Cloud Storage — Google Cloud Platform](https://cloud.google.com/storage/docs/gsutil/commands/defacl#ch)

Like this:

```bash
gsutil cors set example-cors.json gs://example-bucket
gsutil defacl ch -u AllUsers:R gs://example-bucket
```
