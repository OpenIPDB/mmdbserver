package mmdbserver

import (
	"net/http"
	"time"
)

type DownloadHandler func(rw http.ResponseWriter, r *DownloadRequest) error

type DownloadRequest struct {
	EditionID  string
	Date       time.Time
	Suffix     string
	LicenseKey string
}

func (h DownloadHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(rw, ErrMethodNotAllowed)
		return
	}
	var request DownloadRequest
	query := r.URL.Query()
	request.EditionID = query.Get("edition_id")
	request.Suffix = query.Get("suffix")
	request.Date, _ = time.Parse("20060102", query.Get("date"))
	request.LicenseKey = query.Get("license_key")
	var err error
	switch {
	case request.EditionID == "":
		err = ErrInvalidEditionID
	case request.Suffix == "":
		err = ErrInvalidSuffix
	case request.Date.IsZero():
		err = ErrInvalidDate
	case request.LicenseKey == "":
		err = ErrInvalidLicenseKey
	}
	request.Date = request.Date.UTC()
	if err == nil {
		err = h(rw, &request)
	}
	if err != nil {
		WriteError(rw, err)
	}
}
