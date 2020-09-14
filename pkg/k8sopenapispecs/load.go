package k8sopenapispecs

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"regexp"

	"github.com/blang/semver/v4"
	"github.com/gobuffalo/packr/v2"
	openapi_v2 "github.com/googleapis/gnostic/openapiv2"
)

var (
	filenameRegexp = regexp.MustCompile(`k8s_v([\d.]+)_swagger.json.gz`)
	box            = packr.New("k8s-swagger", "./k8s-swagger")
)

// LoadMostRecentVersion loads the openapi spec for the most recent version.
func LoadMostRecentVersion() *openapi_v2.Document {
	var maxVersion semver.Version
	var maxVersionFilename string

	for _, fileName := range box.List() {
		matches := filenameRegexp.FindAllStringSubmatch(fileName, -1)
		if len(matches) == 0 || len(matches[0]) != 2 {
			panic("NOO")
		}
		version := matches[0][1]
		parsedVersion, err := semver.Parse(version)
		if err != nil {
			panic(err)
		}
		if maxVersionFilename == "" || parsedVersion.GT(maxVersion) {
			maxVersion = parsedVersion
			maxVersionFilename = fileName
		}
	}
	if maxVersionFilename == "" {
		panic("NOO")
	}
	contents, err := box.Find(maxVersionFilename)
	if err != nil {
		panic(err)
	}
	gzipR, err := gzip.NewReader(bytes.NewReader(contents))
	if err != nil {
		panic(err)
	}
	decompressedContents, err := ioutil.ReadAll(gzipR)
	if err != nil {
		panic(err)
	}

	parsed, err := openapi_v2.ParseDocument(decompressedContents)
	if err != nil {
		panic(err)
	}
	return parsed
}