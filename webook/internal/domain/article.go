package domain

type Article struct {
	Id         int64
	Title      string
	Content    string
	Author     Author
	Status     ArticleStatus
	CreateTime int64
}

func (a Article) Abstract() string {
	text := []rune(a.Content)
	if len(text) > 100 {
		return string(text[:100])
	}
	return string(text)
}

type ArticleStatus uint8

const (
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnPublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}
func (s ArticleStatus) Valid() bool {
	return s != ArticleStatusUnknown
}

// NonPublished 读者不可见
func (s ArticleStatus) NonPublished() bool {
	return s != ArticleStatusPublished
}

func (s ArticleStatus) String() string {
	switch s {
	case ArticleStatusUnPublished:
		return "unpublished"
	case ArticleStatusPublished:
		return "published"
	case ArticleStatusPrivate:
		return "private"
	default:
		return "unknown"
	}
}

type Author struct {
	Id   int64
	Name string
}
