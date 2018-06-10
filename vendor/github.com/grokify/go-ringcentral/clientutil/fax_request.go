package clientutil

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/gotilla/mime/multipartutil"
	hum "github.com/grokify/gotilla/net/httputilmore"
	"github.com/grokify/gotilla/net/urlutil"
)

const (
	attachmentFieldName = "attachment"
	FaxUrl              = "/restapi/v1.0/account/~/extension/~/fax"
)

func BuildFaxApiUrl(serverUrl string) string {
	return urlutil.JoinAbsolute(serverUrl, FaxUrl)
}

// FaxRequest is a fax request helper that can send more than one attachment.
// The core Swagger Codegen appears to only send one file at a time with a fixed
// MIME part name.
type FaxRequest struct {
	CoverIndex    int
	CoverPageText string
	Resolution    string
	SendTime      *time.Time
	IsoCode       string
	To            []string
	FilePaths     []string
	FileHeaders   []*multipart.FileHeader
}

func NewFaxRequest() FaxRequest {
	return FaxRequest{
		CoverIndex:  -1,
		To:          []string{},
		FilePaths:   []string{},
		FileHeaders: []*multipart.FileHeader{}}
}

func (fax *FaxRequest) builder() (multipartutil.MultipartBuilder, error) {
	builder := multipartutil.NewMultipartBuilder()

	if fax.CoverIndex >= 0 {
		if err := builder.WriteFieldString("coverIndex", strconv.Itoa(fax.CoverIndex)); err != nil {
			return builder, err
		}
	}
	if len(strings.TrimSpace(fax.CoverPageText)) > 0 {
		if err := builder.WriteFieldString("coverPageText", fax.CoverPageText); err != nil {
			return builder, err
		}
	}
	if len(strings.TrimSpace(fax.Resolution)) > 0 {
		if err := builder.WriteFieldString("faxResolution", fax.Resolution); err != nil {
			return builder, err
		}
	}
	if fax.SendTime != nil {
		if err := builder.WriteFieldString("sendTime", fax.SendTime.Format(time.RFC3339)); err != nil {
			return builder, err
		}
	}
	for _, to := range fax.To {
		to := strings.TrimSpace(to)
		if len(to) > 0 {
			if err := builder.WriteFieldString("to", to); err != nil {
				return builder, err
			}
		}
	}

	for _, filePath := range fax.FilePaths {
		if err := builder.WriteFilePath(attachmentFieldName, filePath); err != nil {
			return builder, err
		}
	}
	for _, fileHeader := range fax.FileHeaders {
		if err := builder.WriteFileHeader(attachmentFieldName, fileHeader); err != nil {
			return builder, err
		}
	}

	err := builder.Close()
	return builder, err
}

func (fax *FaxRequest) Post(httpClient *http.Client, url string) (*http.Response, error) {
	builder, err := fax.builder()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, ioutil.NopCloser(builder.Buffer))
	if err != nil {
		return nil, err
	}
	req.Header.Set(hum.HeaderContentType, builder.ContentType())
	return httpClient.Do(req)
}

type FaxCoverPage int

const (
	None         FaxCoverPage = iota // 0
	Ancient                          // 1
	Birthday                         // 2
	Blank                            // 3
	Clasmod                          // 4
	Classic                          // 5
	Confidential                     // 5
	Contempo                         // 7
	Elegant                          // 8
	Express                          // 9
	Formal                           // 10
	Jazzy                            // 11
	Modern                           // 12
	Urgent                           // 13
)

var faxCoverPages = [...]string{
	"None",         // 0
	"Ancient",      // 1
	"Birthday",     // 2
	"Blank",        // 3
	"Clasmod",      // 4
	"Classic",      // 5
	"Confidential", // 6
	"Contempo",     // 7
	"Elegant",      // 8
	"Express",      // 9
	"Formal",       // 10
	"Jazzy",        // 11
	"Modern",       // 12
	"Urgent",       // 13
}

// String returns the English name of the fax cover page ("None", "Ancient", ...).
func (d FaxCoverPage) String() string {
	if None <= d && d <= Urgent {
		return faxCoverPages[d]
	}
	return faxCoverPages[0]
}

func FaxCoverPageNameToIndex(s string) (int, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	l := len(faxCoverPages)
	for i := 0; i < l; i++ {
		name := strings.ToLower(faxCoverPages[i])
		if s == name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("FaxCoverPage NotFound [%v]", s)
}
