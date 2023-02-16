package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/OpenIPDB/mmdbserver"
	"github.com/tg123/go-htpasswd"
)

var (
	databasePath string
	htpasswdPath string
	address      string
	homepage     string
	certFile     string
	keyFile      string
)

func init() {
	flag.StringVar(&databasePath, "db-path", "./databases", "database directory path")
	flag.StringVar(&htpasswdPath, "htpasswd-path", "./htpasswd", "htpasswd file path")
	flag.StringVar(&address, "addr", "localhost:8080", "listen address")
	flag.StringVar(&homepage, "homepage", "https://github.com/OpenIPDB", "homepage url")
	flag.StringVar(&certFile, "cert-file", "", "certificate file")
	flag.StringVar(&keyFile, "key-file", "", "certificate key file")
	flag.Parse()
}

func main() {
	fs := os.DirFS(databasePath)
	passwd, _ := htpasswd.New(htpasswdPath, htpasswd.DefaultSystems, nil)
	handler := &mmdbserver.MMDBUpdateHandler{
		HomePage: homepage,
		MMDBHandler: mmdbserver.MMDBHandlerFunc(func(resp *mmdbserver.Response, r *mmdbserver.Request) (err error) {
			if passwd != nil && !passwd.Match(r.AccountID, r.EditionID) {
				err = mmdbserver.ErrUnauthorized
				return
			}
			database, _ := fs.Open(r.EditionID + ".mmdb")
			_, err = resp.ReadFile(database)
			return
		}),
	}
	var err error
	if certFile != "" && keyFile != "" {
		err = http.ListenAndServeTLS(address, certFile, keyFile, handler)
	} else {
		err = http.ListenAndServe(address, handler)
	}
	log.Println(err)
}
