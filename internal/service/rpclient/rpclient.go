package rpclient

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/net/html"
)

type Client struct {
	BaseUrl string
	Client  http.Client
	Logger  zerolog.Logger
}

func NewClient(baseurl string, timeout time.Duration, log zerolog.Logger) *Client {
	c := http.Client{}
	c.Timeout = timeout
	return &Client{Client: c, BaseUrl: baseurl, Logger: log}

}

type CompanyData struct {
	INN      string
	KPP      string
	Name     string
	HeadName string
}

var ErrCompanyNotFound = errors.New("company not found")

// parse company from https://www.rusprofile.ru/search?query={inn}
func (c *Client) GetCompanyByINN(inn string) (CompanyData, error) {
	baseURL, _ := url.Parse(c.BaseUrl + "/search")
	params := url.Values{}
	params.Add("query", inn)
	params.Add("type", "ul")
	baseURL.RawQuery = params.Encode()

	c.Logger.Debug().Msgf("send request %s", baseURL.String())
	resp, err := c.Client.Get(baseURL.String())
	if err != nil {
		return CompanyData{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return CompanyData{}, ErrCompanyNotFound
	}

	// PARSING <^\_/^>
	c.Logger.Debug().Msgf("GetCompanyByINN/parsing request body...")
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return CompanyData{}, err
	}
	var attrs = []string{
		"clip_inn",
		"clip_kpp",
		"company-name",
		"link-arrow gtm_main_fl"}

	elems := getElementsByAttrs(doc, attrs)
	if len(elems) != 4 {
		return CompanyData{}, ErrCompanyNotFound
	}
	company := CompanyData{}
	for key, node := range elems {
		switch key {
		case "company-name":
			company.Name = node.FirstChild.Data
		case "clip_kpp":
			company.KPP = node.FirstChild.Data
		case "clip_inn":
			company.INN = node.FirstChild.Data
		case "link-arrow gtm_main_fl":
			company.HeadName = node.FirstChild.FirstChild.Data
		}
	}
	c.Logger.Debug().Msgf("GetCompanyByINN/company: %#v", company)
	return company, nil

}
