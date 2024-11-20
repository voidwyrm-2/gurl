package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/akamensky/argparse"
)

func writeFile(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func getHTML(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	}

	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return []byte{}, err
	} else if string(content) == "404: Not Found" {
		return []byte{}, errors.New("404: Not Found")
	}

	return content, nil
}

func main() {
	parser := argparse.NewParser("gurl", "Curl written in Go")

	url := parser.StringList("u", "url", &argparse.Options{Required: true, Default: []string{}, Help: "The URL(s) to download"})
	output := parser.StringList("o", "output", &argparse.Options{Required: false, Default: []string{}, Help: "Write to file instead of stdout"})
	slient := parser.Flag("s", "slient", &argparse.Options{Required: false, Default: false, Help: "Fail sliently"})
	keepGoing := parser.Flag("", "keep-going", &argparse.Options{Required: false, Default: false, Help: "Do not exit the program on download error"})
	useWebname := parser.Flag("", "use-webname", &argparse.Options{Required: false, Default: false, Help: "Use the web name of the downloaded file as the output file name; overrides o/output"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if len(*output) != len(*url) && len(*output) != 0 {
		fmt.Println("the amount of output files must be the same as the amount of URLs")
		os.Exit(1)
	}

	for i, u := range *url {
		if content, err := getHTML(u); err != nil {
			if !*slient {
				fmt.Println(err.Error())
			}
			if !*keepGoing {
				os.Exit(1)
			}
		} else if *useWebname {
			parts := strings.Split(u, "/")
			if err = writeFile(parts[len(parts)-1], content); err != nil {
				if !*slient {
					fmt.Println(err.Error())
				}
				if !*keepGoing {
					os.Exit(1)
				}
			}
		} else if len(*output) != 0 {
			if err = writeFile((*output)[i], content); err != nil {
				if !*slient {
					fmt.Println(err.Error())
				}
				if !*keepGoing {
					os.Exit(1)
				}
			}
		} else {
			fmt.Println(string(content))
		}
	}
}
