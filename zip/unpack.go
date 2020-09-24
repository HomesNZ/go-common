package zip

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// UnpackZIP unpacks the ZIP file at the given path
func Unpack(dir string) (string, error) {
	wd, err := os.Getwd()
	dir = path.Join(wd, dir)

	reader, err := zip.OpenReader(dir)
	if err != nil {
		logrus.WithField("ZIP", "reading the ZIP archive").WithError(err).Error(err)
		return "", err
	}

	unpackPath := strings.Replace(dir, ".zip", "", 1)

	err = os.MkdirAll(unpackPath, 0755)
	if err != nil {
		logrus.WithField("ZIP", "creating ZIP file directory").WithError(err).Error(err)
		return "", err
	}

	for _, file := range reader.File {
		// ensure that any sub folders are created
		path := filepath.Join(unpackPath, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			logrus.WithField("ZIP", "opening a ZIP file").WithError(err).Error(err)
			return "", err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			logrus.WithField("ZIP", "opening ZIP target").WithError(err).Error(err)
			return "", err
		}

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			logrus.WithField("ZIP", "copying ZIP file").WithError(err).Error(err)
			return "", err
		}

		fileReader.Close()
		targetFile.Close()
	}

	return unpackPath, nil
}
