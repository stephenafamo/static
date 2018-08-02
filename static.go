package static

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stephenafamo/mimes"
)

func New(basePath, template_source string) Server {
	server := Server{
		basePath:        basePath,
		template_source: template_source,
	}

	return server
}

type Server struct {
	template_source string
	basePath        string
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var addr strings.Builder

	addr.WriteString(s.basePath)
	addr.WriteString(r.URL.Path)

	file, err := s.GetFile(addr.String(), "")
	checkError(err)

	w.Header().Set("Content-Type", s.getContentType(r.URL.Path))
	w.Write([]byte(file))
}

func (s Server) getContentType(path string) string {
	ext := filepath.Ext(path)
	MIME, _ := mimes.Get(ext)

	var ContentType strings.Builder
	ContentType.WriteString(MIME)
	ContentType.WriteString("; charset=utf-8")

	return ContentType.String()
}

func (s Server) GetFile(name string, extension string) (string, error) {
	// Make a get request
	var addr strings.Builder
	var bodyString string

	addr.WriteString(s.template_source)
	addr.WriteString(name)
	addr.WriteString(extension)

	rs, err := http.Get(addr.String())

	if err != nil {
		return bodyString, err
	}

	if rs.StatusCode < 200 || rs.StatusCode > 299 {
		err = errors.New("Problem getting template--> status code: " + strconv.Itoa(rs.StatusCode))
		return bodyString, err
	}

	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return bodyString, err
	}

	bodyString = string(bodyBytes)

	return bodyString, nil
}

func checkError(err error) {
	if err != nil {
		log.Print(err)
		fmt.Printf("%#v \n", err)
	}
}
