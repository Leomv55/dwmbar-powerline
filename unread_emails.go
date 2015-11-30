package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type gmail struct {
	sync.WaitGroup
	sync.Mutex

	accounts []struct {
		Username, Password string
	}

	client *http.Client

	iteration int
	counts    []int
	results   map[string]int
}

func new_gmail_client(confPath string) (*gmail, error) {
	gm := &gmail{
		client:    &http.Client{},
		iteration: EMAIL_PER_ITERATIONS - 1,
	}

	file, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read gmail account config file: %s - %s", confPath, err)
	}

	if err = json.Unmarshal(file, &gm.accounts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal gmail config file: %s - %s", confPath, err)
	}

	for _ = range gm.accounts {
		gm.counts = append(gm.counts, 0)
	}

	gm.results = make(map[string]int, len(gm.accounts))
	return gm, nil
}

func (gm *gmail) fetch(usr, psw string) (c int, err error) {
	req, err := http.NewRequest("GET", EMAIL_FEED, nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(usr, psw)
	res, err := gm.client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return c, fmt.Errorf(res.Status)
	}

	data := struct {
		Count int `xml:"fullcount"`
	}{}
	err = xml.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return
	}
	return data.Count, nil
}

func unread_emails(confPath string) element {
	gm, err := new_gmail_client(confPath)
	if err != nil {
		return func() (string, error) {
			return "", err
		}
	}

	return func() (string, error) {
		if gm.iteration < EMAIL_PER_ITERATIONS {
			gm.iteration++
			return gm.result(), nil
		}

		gm.Add(len(gm.accounts))
		for _, acc := range gm.accounts {
			go func(u, p string) {
				c, err := gm.fetch(u, p)
				if err != nil {
					log.Println("failed to fetch email count from: %s - %s", u, err)
					c = 0
				}
				gm.Lock()
				gm.results[u] = c
				gm.Unlock()
				gm.Done()
			}(acc.Username, acc.Password)
		}
		gm.Wait()

		var counts []int
		for _, acc := range gm.accounts {
			counts = append(counts, gm.results[acc.Username])
		}
		gm.counts = counts
		gm.iteration = 1
		return gm.result(), nil
	}
}

func (gm *gmail) result() string {
	var out string
	if len(gm.counts) > 0 {
		out = "^i(" + xbm("mail") + ")"
		for _, c := range gm.counts {
			if c > 0 {
				out += fmt.Sprintf(" ^fg(#dc322f)%d^fg()", c)
			} else {
				out += fmt.Sprintf(" %d", c)
			}
		}
	}
	return out
}
