package cookies

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/mengzhuo/cookiestxt"
)

type UserCookie struct {
	cookies []*http.Cookie
}

func NewCookie(uri string) (*UserCookie, error) {
	reader, err := os.ReadFile(uri)
	if err != nil {
		return nil, err
	}
	cookies, err := cookiestxt.Parse(strings.NewReader(string(reader)))
	if err != nil {
		return nil, err
	}
	return &UserCookie{
		cookies: cookies,
	}, nil
}

func (c *UserCookie) AddCookies(request *http.Request) {
	// create a cookie string
	cookieString := ""
	for _, cookie := range c.cookies {
		cookieString += cookie.Name + "=" + cookie.Value + "; "
	}

	// add the cookie string to the request
	request.Header.Add("Cookie", cookieString)
}

func (c *UserCookie) Get(url string) (io.ReadCloser, error) {
	// create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// add the cookies to the request
	addUserDetails(req, c, url)

	// create a new client
	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	var content []byte
	for {
		buf := make([]byte, 1024)
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		if n == 0 {
			break
		}
		content = append(content, buf[:n]...)
	}
	resp.Body.Close()

	// check if under attack
	if strings.Contains(string(content), "Under attack") {
		return nil, fmt.Errorf("under attack")
	}

	return io.NopCloser(strings.NewReader(string(content))), nil
}

func addUserDetails(req *http.Request, uc *UserCookie, url string) {
	// add headers
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en,pt-BR;q=0.8,pt;q=0.5,en-US;q=0.3")
	//req.Header.Add("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Add("Referer", gerReferer(url))
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Sec-Fetch-Dest", "document")
	req.Header.Add("Sec-Fetch-Mode", "navigate")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-User", "?1")
	req.Header.Add("Priority", "u=0, i")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("TE", "trailers")
	uc.AddCookies(req)
}

func gerReferer(target string) string {
	targetURL, err := url.Parse(target)
	if err != nil {
		return target
	}
	return targetURL.Scheme + "://" + targetURL.Host
}
