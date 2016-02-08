package nicome

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type Chat struct {
	Thread    int    `xml:"thread,attr"`
	No        int    `xml:"no,attr"`
	VPos      int    `xml:"vpos,attr"`
	Date      int    `xml:"date,attr"`
	Mail      int    `xml:"mail,attr"`
	UserId    int    `xml:"user_id,attr"`
	Anonymity int    `xml:"anonymity,attr"`
	Text      string `xml:",chardata"`
}

type Packet struct {
	Chat []Chat `xml:"chat"`
}

type client struct {
	mail     string
	password string
	client   *http.Client
}

func NewClient(mail, password string) *client {
	jar, _ := cookiejar.New(nil)
	return &client{mail, password, &http.Client{Jar: jar}}
}

func (c *client) Login() error {
	params := url.Values{
		"mail_tel": []string{c.mail},
		"password": []string{c.password},
	}
	resp, err := c.client.PostForm("https://secure.nicovideo.jp/secure/login?site=niconico", params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *client) Comments(id string, num int) ([]Chat, error) {
	resp, err := c.client.Get("http://www.nicovideo.jp/api/getflv/" + id)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()

	params, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, err
	}

	user_id := ""
	for _, v := range c.client.Jar.Cookies(resp.Request.URL) {
		if strings.HasPrefix(v.Value, "user_session_") {
			user_id = v.Value[13:]
		}
	}
	if user_id == "" {
		return nil, errors.New("can't get user_id")
	}
	if num > 0 {
		num = -num
	}
	payload := fmt.Sprintf(`
				<packet>
			      <thread thread="%s" version="20061206" res_from="%d" user_id="%s"/>
			    </packet>
				`, params.Get("thread_id"), num, user_id)

	resp, err = c.client.Post(params.Get("ms"), "application/x-www-form-urlencoded", strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	var packet Packet
	err = xml.NewDecoder(resp.Body).Decode(&packet)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return packet.Chat, nil
}
