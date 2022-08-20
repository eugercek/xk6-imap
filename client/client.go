package client

import (
	"fmt"
	"io/ioutil"
	"mime/quotedprintable"
	"net/textproto"
	"strconv"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type EmailClient struct {
	Email    string
	Password string
	Url      string
	Port     int
	client   *client.Client
}

func (e *EmailClient) Login() string {
	c, err := client.DialTLS(e.Url+":"+strconv.Itoa(e.Port), nil)

	if err != nil {
		return err.Error()
	}

	e.client = c

	err = e.client.Login(e.Email, e.Password)

	if err != nil {
		return err.Error()
	}

	return ""

}

func (e *EmailClient) Read(header textproto.MIMEHeader) (string, string) {
	_, err := e.client.Select("INBOX", true)

	if err != nil {
		fmt.Println(err)
		return "", err.Error()
	}

	criteria := &imap.SearchCriteria{
		Header: header,
	}

	ids, err := e.client.Search(criteria)

	if err != nil {
		fmt.Println(err)
		return "", err.Error()
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(ids...)

	items := []imap.FetchItem{imap.FetchItem("BODY[TEXT]")}

	messages := make(chan *imap.Message, 1)

	err = e.client.Fetch(seqSet, items, messages)

	if err != nil {
		fmt.Println(err)
		return "", err.Error()
	}

	msg := <-messages

	if msg == nil {
		return "", "No message"
	}

	section, _ := imap.ParseBodySectionName("BODY[TEXT]")
	r := msg.GetBody(section)

	qr := quotedprintable.NewReader(r)
	bs, err := ioutil.ReadAll(qr)

	if err != nil {
		return "", err.Error()
	}

	return string(bs), ""

}

func (e *EmailClient) Logout() {
	e.client.Logout()
}
