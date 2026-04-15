package storage

import "errors"

var (
	// ErrConfigMissing 配置缺失
	ErrConfigMissing = errors.New("storage config is missing")

	// ErrUnsupportedProvider 不支持的存储提供者
	ErrUnsupportedProvider = errors.New("unsupported storage provider")

	// ErrUploadFailed 上传失败
	ErrUploadFailed = errors.New("failed to upload file")

	// ErrDeleteFailed 删除失败
	ErrDeleteFailed = errors.New("failed to delete file")

	// ErrFileNotFound 文件不存在
	ErrFileNotFound = errors.New("file not found")

	// ErrFileTooLarge 文件过大
	ErrFileTooLarge = errors.New("file size exceeds maximum limit")

	// ErrInvalidMimeType 无效的MIME类型
	ErrInvalidMimeType = errors.New("invalid mime type")

	// ErrDownloadFailed 下载失败
	ErrDownloadFailed = errors.New("failed to download file")

	// ErrBucketNotFound 存储桶不存在
	ErrBucketNotFound = errors.New("bucket not found")

	// ErrConnectionFailed 连接失败
	ErrConnectionFailed = errors.New("connection to storage failed")
)
