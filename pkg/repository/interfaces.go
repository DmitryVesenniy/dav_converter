package repository

import "io"

type IDavFile interface {
	GetPathFrame() string
	GetBasePath() string
	GetName() string
	GetReader() io.ReadSeeker
	Close() error
}

type IDavPath interface {
	SetDavPath(string)
	GetDavList(int) ([]IDavFile, error)
	Next() ([]IDavFile, error)
}
