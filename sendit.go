// this Package is from https://github.com/gadelkareem/go-helpers for personal use
// user must pull package into their local repo by using "go get github.com/fl0wmaster6000/sendit"
package sendit

import (
	"crypto/tls"
	"encoding/base64"
	"log"
	"net/smtp"
	"strings"
)

// ex: SendMail("127.0.0.1:25", (&mail.Address{"from name", "from@example.com"}).String(),
// "Email Subject", "message body", []string{(&mail.Address{"to name", "to@example.com"}).String()})

func SendMail(addr string, from string, subject string, body string, to []string) error {
	r := strings.NewReplacer( "\r\n", "", "\r", "", "\n", "", "%0a", "",  "%0d", "")

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         addr,
	}

	c, err := smtp.Dial("tcp", addr, tlsconfig)   // c is a net/smtp Client structure
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if err = c.Mail(r.Replace(from)); err != nil {
		return err
	}

	for i:= range to {
		to[i] = r.Replace(to[i])
		if err = c.Rcpt(to[i]); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	//create the actual smtp mail (not the envelope. Envelope is set above)
	// 'to' is an array of strings

	msg := "To: " + strings.Join(to, ",") + "\r\n" + "From: " + from + "\r\n" + "Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" + "Content-Transfer-Encoding: base64\r\n" + "\r\n" +
		base64.StdEncoding.EncodeToString([]byte(body))  // []byte(body) changes the "body" string to a byte array

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

