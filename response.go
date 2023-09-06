package mmdbserver

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
	"io/fs"
	"net/http"
	"time"
)

type Response struct {
	gzipped bytes.Buffer
	md5Hash hash.Hash
	ModTime time.Time
}

func NewResponse() (response *Response) {
	return &Response{
		md5Hash: md5.New(),
		ModTime: time.Now(),
	}
}

func (r *Response) ReadFile(file fs.File) (n int64, err error) {
	if file == nil {
		err = ErrDatabaseNotFound
		return
	}
	var stat fs.FileInfo
	if stat, err = file.Stat(); err != nil || stat.IsDir() {
		err = ErrDatabaseNotFound
		return
	}
	r.ModTime = stat.ModTime()
	return r.ReadFrom(file)
}

func (r *Response) ReadFrom(file io.Reader) (n int64, err error) {
	if file == nil {
		err = ErrDatabaseNotFound
		return
	}
	r.gzipped.Reset()
	r.md5Hash.Reset()
	return io.Copy(
		io.MultiWriter(gzip.NewWriter(&r.gzipped), r.md5Hash),
		file,
	)
}

func (r *Response) WriteTo(w io.Writer) (n int64, err error) {
	if rw, ok := w.(http.ResponseWriter); ok {
		header := rw.Header()
		header.Set("Content-Type", "application/x-gzip")
		header.Set("Content-Encoding", "gzip")
		header.Set("Last-Modified", r.ModTime.UTC().Format(http.TimeFormat))
		header.Set("X-Database-Md5", hex.EncodeToString(r.MD5Hash()))
	}
	return r.gzipped.WriteTo(w)
}

func (r *Response) MD5Hash() []byte {
	return r.md5Hash.Sum(nil)
}

func (r *Response) Valid() bool {
	return r.gzipped.Len() > 0 && !r.ModTime.IsZero()
}
