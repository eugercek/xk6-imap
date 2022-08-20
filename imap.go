package imap

import (
	"io/ioutil"
	"mime/quotedprintable"
	"net/textproto"
	"strconv"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"

	ec "github.com/eugercek/xk6-imap/client"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/imap", new(Imap))
}

type Imap struct{}

// Simple function for one time read
// Use EmailClient for more complex needs
func (*Imap) Read(email, password, URL string, port int, header textproto.MIMEHeader) (string, string) {
	c, err := client.DialTLS(URL+":"+strconv.Itoa(port), nil)

	if err != nil {
		return "", err.Error()
	}

	defer c.Logout()

	if err := c.Login(email, password); err != nil {
		return "", err.Error()
	}

	_, err = c.Select("INBOX", true)

	if err != nil {
		return "", err.Error()
	}

	criteria := &imap.SearchCriteria{
		Header: header,
	}

	ids, err := c.Search(criteria)

	if err != nil {
		return "", err.Error()
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(ids...)

	items := []imap.FetchItem{imap.FetchItem("BODY[TEXT]")}

	messages := make(chan *imap.Message, 1)

	err = c.Fetch(seqSet, items, messages)

	if err != nil {
		return "", err.Error()
	}

	msg := <-messages

	if msg == nil {
		return "", err.Error()
	}

	section, _ := imap.ParseBodySectionName("BODY[TEXT]")
	r := msg.GetBody(section)

	qr := quotedprintable.NewReader(r)
	bs, err := ioutil.ReadAll(qr)

	if err != nil {
		return "", err.Error()
	}

	return string(bs), "" // TODO Maybe return "OK"
}

// Create new email client
func (*Imap) EmailClient(email, password, url string, port int) *ec.EmailClient {
	return &ec.EmailClient{
		Email:    email,
		Password: password,
		Url:      url,
		Port:     port,
	}
}
