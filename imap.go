package imap

import (
	"fmt"
	"io/ioutil"
	"mime/quotedprintable"
	"net/textproto"
	"strconv"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/imap", new(Imap))
}

type Imap struct{}

func (*Imap) Read(email, password, URL string, port int, header textproto.MIMEHeader) string {
	c, err := client.DialTLS(URL+":"+strconv.Itoa(port), nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer c.Logout()

	if err := c.Login(email, password); err != nil {
		fmt.Println(err)
		return ""
	}

	_, err = c.Select("INBOX", true)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	criteria := &imap.SearchCriteria{
		Header: header,
	}

	ids, err := c.Search(criteria)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(ids...)

	items := []imap.FetchItem{imap.FetchItem("BODY[TEXT]")}

	messages := make(chan *imap.Message, 1)

	err = c.Fetch(seqSet, items, messages)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	msg := <-messages

	if msg == nil {
		fmt.Println("No message")
		return ""
	}

	section, _ := imap.ParseBodySectionName("BODY[TEXT]")
	r := msg.GetBody(section)

	qr := quotedprintable.NewReader(r)
	bs, err := ioutil.ReadAll(qr)

	if err != nil {
		fmt.Println("bs: ", "err is ", err)
	}

	return string(bs)
}
