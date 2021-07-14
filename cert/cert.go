package cert

import (
	"io/ioutil"
	"log"
	"os"
	"path"
)

func init() {
	path := path.Join("cert", "ca-certificate.crt")
	ioutil.WriteFile(path, []byte(os.Getenv("CERT")), 0777)
	log.Println("cert", os.Getenv("CERT"))
}
