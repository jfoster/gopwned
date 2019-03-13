package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

func main() {
	if err := app(); err != nil {
		log.Fatal(err)
	}
}

func app() error {
	for _, arg := range os.Args[1:] {
		sum := sha1.Sum([]byte(arg))
		s := fmt.Sprintf("%x", sum)
		s = strings.ToUpper(s)

		url := fmt.Sprintf("https://api.pwnedpasswords.com/range/%s", s[:5])
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		hashes := strings.Split(string(body), "\r\n")
		occur := 0
		for _, hash := range hashes {
			split := strings.Split(hash, ":")
			if split[0] == s[5:] {
				if in, err := strconv.Atoi(split[1]); err == nil {
					occur = in
				}
			}
		}
		fmt.Printf("%s: %d occurrences found.\n", arg, occur)
	}
	return nil
}
