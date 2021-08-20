package zerogame

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func currentPlatform() string {
	return runtime.GOOS
}

func getURL(u string) ([]byte, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	switch parsedURL.Scheme {
	case "file":
		return ioutil.ReadFile(u[7:])
	case "http", "https":
		client := &http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		res, err := client.Get(u)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		return ioutil.ReadAll(res.Body)
	default:
		return nil, fmt.Errorf("invalid scheme: %q. must be one of: [file http https] ", u)
	}
}

func removeFileExtension(filename string) string {
	return filename[0 : len(filename)-len(filepath.Ext(filename))]
}

func replaceFileExtension(filename, extension string) string {
	return filename[0:len(filename)-len(filepath.Ext(filename))] + "." + extension
}

func extract(path string) (files []string, err error) {
	if filepath.Ext(path) != ".zip" {
		return nil, errors.New("only zip archives are supported")
	}
	return unzip(path, removeFileExtension(path))
}

func unzip(src, dest string) ([]string, error) {
	var filenames []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
