# Pulumi GCP Go: Static Website

This project provisions a GCP Cloud Storage bucket for static website hosting using Pulumi and Go. It demonstrates how to:
- Use the Pulumi GCP provider in a Go program
- Create a storage bucket configured for website hosting
- Generate unique bucket names with random suffixes
- Configure public access for website content

## Providers

- Google Cloud Platform via the Pulumi GCP SDK for Go (`github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp`)
- Random via the Pulumi Random SDK for Go (`github.com/pulumi/pulumi-random/sdk/v4/go/random`)

## Resources

- **Random ID** (`random.RandomId`)
  - Generates a random suffix for bucket name uniqueness

- **Storage Bucket** (`gcp.storage.Bucket`)
  - Bucket configured for static website hosting
  - STANDARD storage class for Free Tier eligibility
  - Configurable index and 404 pages

- **Bucket IAM Binding** (`gcp.storage.BucketIAMBinding`)
  - Grants public read access to allUsers (optional)

## Outputs

- **bucketName**: The name of the bucket
- **bucketUrl**: The base URL of the bucket
- **websiteEndpoint**: GCS virtual-hosted endpoint for the bucket
- **objectBaseUrl**: The blob base URL for individual objects
- **bucketSelfLink**: The URI of the created resource

## When to Use This Project

- You want to host a static website on GCP for free/low cost
- You need a simple, serverless hosting solution
- You want to use a custom domain with Cloud CDN

## Prerequisites

- Go 1.21+ installed
- A Google Cloud account with billing enabled
- GCP credentials configured for Pulumi (via `gcloud auth application-default login`)
- Cloud Storage API enabled

## Usage

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Configure your stack:
   ```bash
   cp Pulumi.dev.yaml.example Pulumi.dev.yaml
   pulumi config set gcp:project YOUR_PROJECT_ID
   pulumi config set website:bucketName my-static-site
   ```

3. Preview and deploy:
   ```bash
   pulumi preview
   pulumi up
   ```

4. Upload website files:
   ```bash
   gsutil -m cp -r website/* gs://<bucket-name>/
   gsutil -m acl ch -r -u AllUsers:R gs://<bucket-name>/
   ```

## Project Layout

```
├── Pulumi.yaml                  Pulumi project definition
├── Pulumi.dev.yaml.example      Template for local dev configuration
├── go.mod                       Go module declaration and dependencies
├── main.go                      Pulumi program defining website resources
├── .gitignore                   Git ignore rules
└── LICENSE                      MIT License
```

## Configuration

| Name | Description | Default |
|------|-------------|---------|
| `gcp:project` | The Google Cloud project to deploy into | _required_ |
| `gcp:region` | The GCP region for resources | `us-central1` |
| `website:bucketName` | Base name of the bucket (required) | _required_ |
| `website:mainPageSuffix` | Index page for directories | `index.html` |
| `website:notFoundPage` | 404 page | `404.html` |
| `website:enablePublicAccess` | Grant public read access | `true` |
| `website:websiteContentPath` | Path to website files | `website` |

## Next Steps

- Add custom domain mapping with Cloud CDN
- Configure Cache-Control headers for assets
- Set up CI/CD for automatic deployments
- Add error page customization
- Enable bucket logging for analytics

## Getting Help

- Pulumi Documentation: https://www.pulumi.com/docs/
- GCP Provider Reference: https://www.pulumi.com/registry/packages/gcp/
- Community Slack: https://slack.pulumi.com/
- GitHub Issues: https://github.com/pulumi/pulumi/issues