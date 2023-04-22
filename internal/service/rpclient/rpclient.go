package rpclient

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

type Client struct {
	BaseUrl string
	Client  http.Client
}

func NewClient(baseurl string, timeout time.Duration) *Client {
	c := http.Client{}
	c.Timeout = timeout
	return &Client{Client: c}

}

type CompanyData struct {
	INN      string
	KPP      string
	Name     string
	HeadName string
}

func GetAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

func checkId(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := GetAttribute(n, "id")
		if ok && s == id {
			return true
		}
	}
	return false
}

func traverse(n *html.Node, id string) *html.Node {
	if checkId(n, id) {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverse(c, id)
		if result != nil {
			return result
		}
	}

	return nil
}

func getElementById(n *html.Node, id string) *html.Node {
	return traverse(n, id)
}

// id="clip_inn"
// id="clip_kpp"
// class="company-name"
// class="link-arrow gtm_main_fl -> span   ФИО
const (
	parse_inn_id   = "clip_inn"
	parse_kpp_id   = "clip_kpp"
	parse_name     = "company_name"
	parse_headname = "link-arrow gtm_main_fl"
)

// parse company from https://www.rusprofile.ru/search?query=500100732259
func (c *Client) GetCompanyByINN(inn string) (CompanyData, error) {
	baseURL, _ := url.Parse(c.BaseUrl)
	params := url.Values{}
	params.Add("query", inn)
	baseURL.RawQuery = params.Encode()

	resp, err := c.Client.Get(baseURL.String())
	if err != nil {
		return CompanyData{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return CompanyData{}, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	// PARSING
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return CompanyData{}, err
	}

	inn_ := getElementById(doc, parse_inn_id).Data
	fmt.Println(inn_)

	return CompanyData{}, nil

}
