package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/storage"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		projectID := config.Require(ctx, "gcp:project")
		region := config.Get(ctx, "gcp:region")
		if region == "" {
			region = "us-central1"
		}

		bucketName := config.Require(ctx, "website:bucketName")
		mainPageSuffix := config.Get(ctx, "website:mainPageSuffix")
		if mainPageSuffix == "" {
			mainPageSuffix = "index.html"
		}
		notFoundPage := config.Get(ctx, "website:notFoundPage")
		if notFoundPage == "" {
			notFoundPage = "404.html"
		}
		enablePublicAccess := config.GetBool(ctx, "website:enablePublicAccess")
		websiteContentPath := config.Get(ctx, "website:websiteContentPath")
		if websiteContentPath == "" {
			websiteContentPath = "website"
		}

		// Generate random suffix for bucket name
		bucketSuffix, err := random.NewRandomId(ctx, "bucket-suffix", &random.RandomIdArgs{
			ByteLength: pulumi.Int(4),
		})
		if err != nil {
			return err
		}

		// Create storage bucket for website
		websiteBucket, err := storage.NewBucket(ctx, "website-bucket", &storage.BucketArgs{
			Name:                     pulumi.Sprintf("%s-%s", bucketName, bucketSuffix.Hex),
			Location:                 pulumi.String(region),
			StorageClass:             pulumi.String("STANDARD"),
			ForceDestroy:             pulumi.Bool(true),
			UniformBucketLevelAccess: pulumi.Bool(true),
			Website: &storage.BucketWebsiteArgs{
				MainPageSuffix: pulumi.String(mainPageSuffix),
				NotFoundPage:   pulumi.String(notFoundPage),
			},
		})
		if err != nil {
			return err
		}

		// Grant public access if enabled
		if enablePublicAccess {
			_, err = storage.NewBucketIAMBinding(ctx, "public-access", &storage.BucketIAMBindingArgs{
				Bucket: websiteBucket.Name,
				Role:   pulumi.String("roles/storage.objectViewer"),
				Members: pulumi.StringArray{
					pulumi.String("allUsers"),
				},
			})
			if err != nil {
				return err
			}
		}

		// Note: In Go, we can't easily walk the filesystem at deploy time
		// This would typically be done via a separate upload step or CI/CD
		// For now, we'll export the bucket info for manual upload

		// Export outputs
		ctx.Export("projectID", pulumi.String(projectID))
		ctx.Export("websiteContentPath", pulumi.String(websiteContentPath))
		ctx.Export("bucketName", websiteBucket.Name)
		ctx.Export("bucketUrl", websiteBucket.Url)
		ctx.Export("websiteEndpoint", pulumi.Sprintf("https://%s.storage.googleapis.com/", websiteBucket.Name))
		ctx.Export("objectBaseUrl", pulumi.Sprintf("https://storage.googleapis.com/%s/", websiteBucket.Name))
		ctx.Export("bucketSelfLink", websiteBucket.SelfLink)

		return nil
	})
}
