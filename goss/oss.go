package goss

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Client struct {
	Svc *s3.S3
}

func NewS3Client(cfg *aws.Config) *S3Client {
	client := &S3Client{}
	sess := session.Must(session.NewSession())
	client.Svc = s3.New(sess, cfg)
	return client
}

func NewS3ClientWithSharedCredentials(endpoint, filename, profile string) *S3Client {
	// 这个区域的配置暂时是随意填写的，不给为空
	region := "ml"
	// 强制使用路径区分桶，配置默认用子域名区分桶
	forcePathStyle := true
	config := &aws.Config{
		Credentials:      credentials.NewSharedCredentials(filename, profile),
		Endpoint:         &endpoint,
		Region:           &region,
		S3ForcePathStyle: &forcePathStyle,
	}
	return NewS3Client(config)
}

func NewS3ClientWithStaticCredentials(endpoint, access_key_id, access_secret_key, token string) *S3Client {
	// 这个区域的配置暂时是随意填写的，不给为空
	region := "ml"
	// 强制使用路径区分桶，配置默认用子域名区分桶
	forcePathStyle := true
	config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(access_key_id, access_secret_key, token),
		Endpoint:         &endpoint,
		Region:           &region,
		S3ForcePathStyle: &forcePathStyle,
	}

	//if glog.IsDebug() {
	//	config.LogLevel = aws.LogLevel(aws.LogDebug)
	//}

	return NewS3Client(config)
}

func (this *S3Client) GetObject(bucket string, key string) (*s3.GetObjectOutput, error) {
	result, err := this.Svc.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	return result, err
}

func (this *S3Client) PutObject(bucket string, key string, input io.Reader, metadata map[string]string) error {

	objectInput := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(input),
		Bucket: &bucket,
		Key:    &key,
	}

	m := make(map[string]*string)
	for key, value := range metadata {
		m[key] = &value
	}
	objectInput.SetMetadata(m)

	_, err := this.Svc.PutObject(objectInput)
	return err
}

func (this *S3Client) GetObjectMetadata(bucket string, key string) (map[string]string, error) {
	var tagging = make(map[string]string)
	result, err := this.Svc.HeadObject(&s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		return tagging, err
	}

	for key, value := range result.Metadata {
		tagging[key] = *value
	}
	return tagging, err
}

func (this *S3Client) GetObjectTagging(bucket string, key string) (map[string]string, error) {
	var tagging = make(map[string]string)
	result, err := this.Svc.GetObjectTagging(&s3.GetObjectTaggingInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		return tagging, err
	}

	for _, item := range result.TagSet {
		tagging[*item.Key] = *item.Value
	}
	return tagging, err
}

func (this *S3Client) PutObjectTagging(bucket string, key string, tagging map[string]string) error {

	t := &s3.Tagging{}
	for key, value := range tagging {
		t.TagSet = append(t.TagSet, &s3.Tag{
			Key:   &key,
			Value: &value,
		})
	}

	_, err := this.Svc.PutObjectTagging(&s3.PutObjectTaggingInput{
		Bucket:  &bucket,
		Key:     &key,
		Tagging: t,
	})
	return err
}

func (this *S3Client) DeleteObject(bucket string, key string) error {

	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	_, err := this.Svc.DeleteObject(deleteObjectInput)

	return err
}

func (this *S3Client) CopyObject(bucket, source, target string) error {
	input := &s3.CopyObjectInput{
		Bucket:     &bucket,
		CopySource: &source,
		Key:        &target,
	}

	_, err := this.Svc.CopyObject(input)
	return err
}

func (this *S3Client) CopyObjectInSameBucket(bucket, source, target string) error {
	return this.CopyObject(bucket, bucket+"/"+source, target)
}

func (this *S3Client) ListBucket() (*s3.ListBucketsOutput, error) {
	input := &s3.ListBucketsInput{}
	return this.Svc.ListBuckets(input)
}
