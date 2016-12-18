package types

import "time"

type Blob struct {
	Length	int64
	Content []byte
	Name 	string
}

type FileID string

type File struct {
	ID		FileID
	Name 		string
	Length 		int64
	Createtime 	time.Time
	Status 		FileStatus
}

type FileStatus int

var (
	FileInQueue = FileStatus(0)
	FileWriting = FileStatus(1)
	FileSync = FileStatus(2)
	FileOK = FileStatus(3)
	FileDelete = FileStatus(4)
)

func NewFile(name string, length int64) *File {
	return &File{
		Name: name,
		Length: length,
		Createtime: time.Now(),
		Status: FileInQueue,
	}
}