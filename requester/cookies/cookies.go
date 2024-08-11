package cookies

import (
	"net/http"
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
