package m

type GCSBucket struct {
	Name   string
	Region string
}

type GCSModel struct {
}

type GCSManager interface {
	Bucket() string
	SetBucket(bucket string) error
	AvailableBuckets() []GCSBucket
	Prefix() string
	setPrefix(prefix string) error
	MoveUp() error
	MoveDown(prefix string) error
	List() (prefixes []string, keys []string, err error)
	ListObjects(key string) []interface{}
	Download(object interface{}, destPath string) (n int64, err error)
	upperPrefix(prefix string) string
	lastPartPrefix(prefix string) string
}

func NewGCSManager() GCSManager {
	return &GCSModel{}
}

func (gcsm *GCSModel) Bucket() string {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) SetBucket(bucket string) error {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) AvailableBuckets() []GCSBucket {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) Prefix() string {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) setPrefix(prefix string) error {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) MoveUp() error {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) MoveDown(prefix string) error {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) List() (prefixes []string, keys []string, err error) {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) ListObjects(key string) []interface{} {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) Download(object interface{}, destPath string) (n int64, err error) {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) upperPrefix(prefix string) string {
	panic("not implemented") // TODO: Implement
}

func (gcsm *GCSModel) lastPartPrefix(prefix string) string {
	panic("not implemented") // TODO: Implement
}
