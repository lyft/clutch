package github

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"

	githubv3 "github.com/google/go-github/v54/github"
)

const (
	// From https://docs.github.com/en/rest/using-the-rest-api/getting-started-with-the-rest-api?apiVersion=2022-11-28#media-types
	// and https://docs2.lfe.io/v3/media/
	AcceptGithubRawMediaType = "application/vnd.github.v3.raw+json"
	// From https://cs.opensource.google/go/go/+/master:src/encoding/json/scanner.go;l=593?q=scanner.go&ss=go%2Fgo
	JSONSyntaxErrorPrefix = "invalid character"
)

// RegEx only for `Repositories.GetContents` must match something like `/repos/lyft/clutch/contents/README.md`
var GetRepositoryContentRegex = regexp.MustCompile(`^\/repos\/[\w-]+\/[\w-]+\/contents\/[\w-.\/]+$`)

// `Repositories.GetContents“ method cannot process a raw response, so we intercept it
func InterceptGetRepositoryContentResponse(res *http.Response) error {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var fileContentBytes []byte
	var fileContent *githubv3.RepositoryContent

	// Sometimes we get a githubv3.RepositoryContent body and we have to check
	err = json.Unmarshal(body, &fileContent)
	switch {
	case fileContent != nil &&
		fileContent.Content != nil &&
		fileContent.Encoding != nil &&
		fileContent.Path != nil &&
		fileContent.GitURL != nil:
		fileContentBytes = body
	case err != nil:
		if !strings.HasPrefix(err.Error(), JSONSyntaxErrorPrefix) {
			return err
		}
		fallthrough
	default:
		// Set the body as a `githubv3.RepositoryContent`
		fileContentBytes, err = json.Marshal(&githubv3.RepositoryContent{
			Content:  githubv3.String(string(body)),
			Encoding: githubv3.String(""), // The content is not encoded
		})
		if err != nil {
			return err
		}
	}

	// Recreate the response body
	res.Body = io.NopCloser(bytes.NewReader(fileContentBytes))
	res.ContentLength = int64(len(fileContentBytes))

	return nil
}
