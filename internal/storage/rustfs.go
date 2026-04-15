package storage

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gogf/gf/v2/frame/g"
)

// RustFSProvider RustFS 存储提供者 (S3 兼容)
type RustFSProvider struct {
	Client     *s3.Client
	Uploader   *manager.Uploader
	Bucket     string
	Endpoint   string
	Timeout    time.Duration
	PathPrefix string
}

// NewRustFSProvider 创建 RustFS 存储提供者
func NewRustFSProvider(cfg *RustFSConfig) (*RustFSProvider, error) {
	if cfg == nil {
		return nil, ErrConfigMissing
	}

	// 设置默认超时
	timeout := time.Duration(cfg.Timeout) * time.Second
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	// 创建 AWS 配置
	// RustFS 使用 path-style URL，需要禁用虚拟主机样式
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               cfg.Endpoint,
			SigningRegion:     "us-east-1", // RustFS 不关心区域
			HostnameImmutable: true,        // 保持主机名不变
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.Username,
			cfg.Password,
			"",
		)),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// 创建 S3 客户端
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true // 使用 path-style URL
	})

	// 创建上传器
	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // 5MB 分片
	})

	return &RustFSProvider{
		Client:     client,
		Uploader:   uploader,
		Bucket:     cfg.Bucket,
		Endpoint:   cfg.Endpoint,
		Timeout:    timeout,
		PathPrefix: "",
	}, nil
}

// SetPathPrefix 设置路径前缀
func (p *RustFSProvider) SetPathPrefix(prefix string) {
	p.PathPrefix = prefix
}

// Upload 上传文件到 RustFS
func (p *RustFSProvider) Upload(ctx context.Context, reader io.Reader, objectKey string, contentType string) (string, error) {
	// 构建完整的对象键
	fullKey := objectKey
	if p.PathPrefix != "" {
		fullKey = fmt.Sprintf("%s/%s", strings.TrimSuffix(p.PathPrefix, "/"), objectKey)
	}

	g.Log().Infof(ctx, "Uploading file to RustFS: bucket=%s, key=%s, contentType=%s", p.Bucket, fullKey, contentType)

	// 设置上传超时
	ctx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	// 执行上传
	_, uploadErr := p.Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(p.Bucket),
		Key:         aws.String(fullKey),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if uploadErr != nil {
		g.Log().Errorf(ctx, "Failed to upload file to RustFS: %v", uploadErr)
		return "", fmt.Errorf("%w: %v", ErrUploadFailed, uploadErr)
	}

	// 返回相对路径（不含 endpoint），由前端通过图片代理接口访问
	g.Log().Infof(ctx, "File uploaded successfully: %s", fullKey)
	return fullKey, nil
}

// Delete 从 RustFS 删除文件
func (p *RustFSProvider) Delete(ctx context.Context, objectKey string) error {
	// 构建完整的对象键
	fullKey := objectKey
	if p.PathPrefix != "" {
		fullKey = fmt.Sprintf("%s/%s", strings.TrimSuffix(p.PathPrefix, "/"), objectKey)
	}

	g.Log().Infof(ctx, "Deleting file from RustFS: bucket=%s, key=%s", p.Bucket, fullKey)

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	// 执行删除
	_, err := p.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(fullKey),
	})
	if err != nil {
		g.Log().Errorf(ctx, "Failed to delete file from RustFS: %v", err)
		return fmt.Errorf("%w: %v", ErrDeleteFailed, err)
	}

	g.Log().Infof(ctx, "File deleted successfully: %s", fullKey)
	return nil
}

// Download 从 RustFS 下载文件
func (p *RustFSProvider) Download(ctx context.Context, objectKey string) ([]byte, error) {
	// 构建完整的对象键
	fullKey := objectKey
	if p.PathPrefix != "" {
		fullKey = fmt.Sprintf("%s/%s", strings.TrimSuffix(p.PathPrefix, "/"), objectKey)
	}

	g.Log().Infof(ctx, "Downloading file from RustFS: bucket=%s, key=%s", p.Bucket, fullKey)

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	// 执行下载
	result, err := p.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(fullKey),
	})
	if err != nil {
		g.Log().Errorf(ctx, "Failed to download file from RustFS: %v", err)
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer result.Body.Close()

	// 读取文件内容
	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	g.Log().Infof(ctx, "File downloaded successfully: %s (%d bytes)", fullKey, len(data))
	return data, nil
}

// Copy 在 RustFS 内复制文件
func (p *RustFSProvider) Copy(ctx context.Context, srcKey string, dstKey string) error {
	// 构建完整的对象键
	fullSrcKey := srcKey
	fullDstKey := dstKey
	if p.PathPrefix != "" {
		prefix := strings.TrimSuffix(p.PathPrefix, "/")
		fullSrcKey = fmt.Sprintf("%s/%s", prefix, srcKey)
		fullDstKey = fmt.Sprintf("%s/%s", prefix, dstKey)
	}

	g.Log().Infof(ctx, "Copying file in RustFS: %s -> %s", fullSrcKey, fullDstKey)

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	// 执行复制 (S3 CopyObject)
	copySource := fmt.Sprintf("%s/%s", p.Bucket, fullSrcKey)
	_, err := p.Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(p.Bucket),
		Key:        aws.String(fullDstKey),
		CopySource: aws.String(copySource),
	})
	if err != nil {
		g.Log().Errorf(ctx, "Failed to copy file in RustFS: %v", err)
		return fmt.Errorf("failed to copy file: %w", err)
	}

	g.Log().Infof(ctx, "File copied successfully: %s -> %s", fullSrcKey, fullDstKey)
	return nil
}

// GetURL 获取文件访问 URL
func (p *RustFSProvider) GetURL(objectKey string) string {
	// 构建完整的对象键
	fullKey := objectKey
	if p.PathPrefix != "" {
		fullKey = fmt.Sprintf("%s/%s", strings.TrimSuffix(p.PathPrefix, "/"), objectKey)
	}

	// 返回直接访问 URL (path-style)
	return fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(p.Endpoint, "/"), p.Bucket, fullKey)
}

// HealthCheck 健康检查
func (p *RustFSProvider) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 尝试列出存储桶内容来验证连接
	_, err := p.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(p.Bucket),
		MaxKeys: aws.Int32(1),
	})
	if err != nil {
		return fmt.Errorf("%w: %v", ErrConnectionFailed, err)
	}

	return nil
}

// PresignURL 生成预签名 URL (可选功能)
func (p *RustFSProvider) PresignURL(ctx context.Context, objectKey string, expireDuration time.Duration) (string, error) {
	// 构建完整的对象键
	fullKey := objectKey
	if p.PathPrefix != "" {
		fullKey = fmt.Sprintf("%s/%s", strings.TrimSuffix(p.PathPrefix, "/"), objectKey)
	}

	presignClient := s3.NewPresignClient(p.Client)
	presignedReq, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.Bucket),
		Key:    aws.String(fullKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expireDuration
	})
	if err != nil {
		return "", fmt.Errorf("failed to presign URL: %w", err)
	}

	return presignedReq.URL, nil
}

// CheckBucketExists 检查存储桶是否存在
func (p *RustFSProvider) CheckBucketExists(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := p.Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(p.Bucket),
	})
	if err != nil {
		// 检查是否是桶不存在的错误
		errStr := err.Error()
		if strings.Contains(errStr, "NotFound") || strings.Contains(errStr, "404") || strings.Contains(errStr, "NoSuchBucket") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// CreateBucket 创建存储桶
func (p *RustFSProvider) CreateBucket(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	_, err := p.Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(p.Bucket),
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	return nil
}
