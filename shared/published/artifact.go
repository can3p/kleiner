package published

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/can3p/kleiner/shared/types"
	"github.com/go-resty/resty/v2"
	"github.com/minio/selfupdate"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func compileDownloadURL(buildinfo *types.BuildInfo, tag string) string {
	return fmt.Sprintf(
		"https://github.com/%s/releases/download/v%s/%s",
		buildinfo.GithubRepo, tag, compileArtifactName(buildinfo))
}

func compileArtifactName(buildinfo *types.BuildInfo) string {
	arch := buildinfo.Architecture

	switch arch {
	case "amd64":
		arch = "x86_64"
	case "386":
		arch = "i386"
	case "arm68":
		arch = "armv64"
	}

	return fmt.Sprintf("%s_%s_%s.tar.gz", buildinfo.ProjectName, cases.Title(language.English).String(runtime.GOOS), arch)
}

func DownloadAndReplaceBinary(buildinfo *types.BuildInfo, tag string) error {
	b, err := downloadBinary(buildinfo, tag)

	if err != nil {
		return err
	}

	err = selfupdate.Apply(bytes.NewReader(b), selfupdate.Options{})
	if err != nil {
		return errors.Wrapf(err, "failed to update the binary")
	}

	return nil
}

func downloadBinary(buildinfo *types.BuildInfo, tag string) ([]byte, error) {
	url := compileDownloadURL(buildinfo, tag)
	fileToExtract := buildinfo.ProjectName

	client := resty.New()

	resp, err := client.R().Get(url)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to download the binary, url = %s", url)
	}

	body := resp.Body()

	if rc := resp.StatusCode(); rc != http.StatusOK {
		return nil, errors.Errorf("failed to download the binary, url = %s, status code = %d, body = [%s]", url, rc, string(body))
	}

	bodyReader := bytes.NewReader(body)
	gzipReader, err := gzip.NewReader(bodyReader)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to decompress gzip, url = %s, size = %d bytes", url, len(body))
	}

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if header.Typeflag == tar.TypeReg && header.Name == fileToExtract {
			b, err := io.ReadAll(tarReader)

			if err != nil {
				return nil, err
			}
			return b, nil
		}
	}

	return nil, errors.Errorf("file with the name [%s] has not been found in the archive", fileToExtract)
}
