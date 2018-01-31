package wechat

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/labstack/echo"
)

const (
	IPAY_WECHAT_PREPAY       = "IPAY_WECHAT_PREPAY"
	IPAY_WECHAT_PREPAY_INNER = "IPAY_WECHAT_PREPAY_INNER"
	IPAY_WECHAT_PREPAY_ERROR = "IPAY_WECHAT_PREPAY_ERROR"
)

func SetCookie(key, value string, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = key
	value = url.QueryEscape(value)
	cookie.Value = value
	cookie.Domain = "p2shop.cn"
	cookie.Path = "/"
	//cookie.Expires = time.Now().Add(1 * time.Hour)
	c.SetCookie(cookie)
}

func SetCookieObj(key string, value interface{}, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = key
	b, _ := json.Marshal(value)
	cookie.Value = url.QueryEscape(string(b))
	cookie.Domain = "p2shop.cn"
	cookie.Path = "/"
	//cookie.Expires = time.Now().Add(1 * time.Hour)
	c.SetCookie(cookie)
}

func ChinaDatetime() (date time.Time) {
	date = time.Now().UTC().Add(8 * time.Hour)
	return
}

func POSTXml(token, url, param string, v interface{}) (resp *http.Response, err error) {
	b := []byte(param)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		err = fmt.Errorf("HTTP New Request Error: %s", err)
		return
	}
	httpReq.Header.Set("Content-Type", "application/xml")
	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err = (&http.Client{}).Do(httpReq)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("[%d %s]%s", resp.StatusCode, resp.Status, string(b))
		return
	}
	if v != nil {
		dec := xml.NewDecoder(resp.Body)
		if err = dec.Decode(&v); err != nil {
			return
		}
	}

	return
}
