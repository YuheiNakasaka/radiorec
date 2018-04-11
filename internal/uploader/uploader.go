package uploader

// Uploader is interface of uploading tasks
type Uploader interface {
	Upload(string, string) error
}
