package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"log"
)

type client interface {
	Get(method string, res interface{}) error
}

type Client struct {
	protocol string
	hostname string
	port     string

	client *http.Client
}

func NewClient(protocol string, hostname string, port string) client {
	client := http.Client{Timeout: time.Duration(30 * time.Second)}

	return Client{
		protocol: protocol,
		hostname: hostname,
		port:     port,
		client:   &client,
	}
}

func (g Client) buildURL(method string) (string, error) {
	u := fmt.Sprintf("%s://%s:%s", g.protocol, g.hostname, g.port)
	url, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	url.Path = path.Join(url.Path, method)
	return url.String(), nil
}

func (g Client) Get(method string, res interface{}) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	url, err := g.buildURL(method)
	if err != nil {
		log.Println(err)
		return err
	}

	r, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return err
	}
	defer r.Body.Close()

	log.Printf("Client GET - %s - %d", url, r.StatusCode)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	if r.StatusCode != 200 {
		return fmt.Errorf("not valid response")
	}

	err = json.Unmarshal(b, res)
	if err != nil {
		return err
	}

	return nil
}
