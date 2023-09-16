package mmdbserver

import (
	"bytes"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
)

type Request struct {
	AccountID  string `json:"account_id"`
	LicenseKey string `json:"license_key"`
	EditionID  string `json:"edition_id"`
	MD5Hash    []byte `json:"md5_hash"`
}

func (r *Request) NewDeploy() bool {
	var empty [16]byte
	return r.EqualHash(empty[:])
}

func (r *Request) EqualHash(md5Hash []byte) bool {
	return bytes.EqualFold(r.MD5Hash, md5Hash)
}

func NewRequest(r *http.Request) (request *Request, err error) {
	request = new(Request)
	request.MD5Hash, _ = hex.DecodeString(r.URL.Query().Get("db_md5"))
	request.EditionID, _ = getEditionId(r.URL.Path)
	request.AccountID, request.LicenseKey, _ = r.BasicAuth()
	switch {
	case request.EditionID == "":
		err = ErrInvalidEditionID
	case !isNumeric(request.AccountID):
		err = ErrInvalidAccountID
	case request.LicenseKey == "":
		err = ErrInvalidAccountID
	case len(request.MD5Hash) != 16:
		err = ErrInvalidMD5Hash
	}
	return
}

func isNumeric(input string) bool {
	_, err := strconv.Atoi(input)
	return err == nil
}

func getEditionId(pathname string) (editionId string, err error) {
	matched := reUpdatePath.FindStringSubmatch(pathname)
	if matched == nil {
		return
	}
	index := reUpdatePath.SubexpIndex("edition")
	return url.PathUnescape(matched[index])
}
