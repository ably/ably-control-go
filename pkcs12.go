package control

import (
	"bytes"
	"fmt"
	"mime/multipart"

	"github.com/hashicorp/go-retryablehttp"
	"software.sslmate.com/src/go-pkcs12"
)

// UploadPKCS12 uploads a PKCS12 certificate file for the specified application.
// The p12File parameter is the raw bytes of the .p12 file, and p12Pass is the password.
func (c *Client) UploadPKCS12(appID string, p12File []byte, p12Pass string) (App, error) {
	_, cert, err := pkcs12.Decode(p12File, p12Pass)
	if err != nil {
		return App{}, fmt.Errorf("invalid PKCS12 file: %w", err)
	}
	if !cert.BasicConstraintsValid {
		return App{}, fmt.Errorf("invalid PKCS12 certificate: missing BasicConstraints extension")
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("p12File", "certificate.p12")
	if err != nil {
		return App{}, err
	}
	_, err = part.Write(p12File)
	if err != nil {
		return App{}, err
	}

	err = writer.WriteField("p12Pass", p12Pass)
	if err != nil {
		return App{}, err
	}

	err = writer.Close()
	if err != nil {
		return App{}, err
	}

	path := fmt.Sprintf("/apps/%s/pkcs12", appID)
	req, err := retryablehttp.NewRequest("POST", c.Url+path, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return App{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var out App
	err = c.requestRaw(req, path, &out)
	return out, err
}
