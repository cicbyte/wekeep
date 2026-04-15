package model

// ImagesInfo 图片信息
type ImagesInfo struct {
	Id             uint   `orm:"id" json:"id"`
	OriginalUrl    string `orm:"original_url" json:"originalUrl"`
	OriginalUrlHash string `orm:"original_url_hash" json:"originalUrlHash"`
	StoragePath    string `orm:"storage_path" json:"storagePath"`
	StorageUrl     string `orm:"storage_url" json:"storageUrl"`
	FileSize       int    `orm:"file_size" json:"fileSize"`
	MimeType       string `orm:"mime_type" json:"mimeType"`
	RefCount       int    `orm:"ref_count" json:"refCount"`
	DownloadStatus int    `orm:"download_status" json:"downloadStatus"`
	ErrorMessage   string `orm:"error_message" json:"errorMessage"`
}

// ArticleImagesInfo 文章-图片关联信息
type ArticleImagesInfo struct {
	Id        uint `orm:"id" json:"id"`
	ArticleId uint `orm:"article_id" json:"articleId"`
	ImageId   uint `orm:"image_id" json:"imageId"`
	Position  int  `orm:"position" json:"position"`
}

// DownloadStatus 下载状态常量
const (
	DownloadStatusPending   = 0 // 待下载
	DownloadStatusDownloading = 1 // 下载中
	DownloadStatusSuccess   = 2 // 下载成功
	DownloadStatusFailed    = 3 // 下载失败
)

// ImageProcessResult 图片处理结果
type ImageProcessResult struct {
	OriginalURL string // 原始URL
	StorageURL  string // 存储后的URL
	ImageID     uint   // 图片ID
	IsNew       bool   // 是否新创建
	Error       error  // 错误信息
}
