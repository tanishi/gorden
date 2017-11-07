package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

const (
	ExitCodeOK = iota
	ExitCodeParseFlagError
)

const (
	envUserID = "GIT_COMMITTER_NAME"
)

type CLI struct {
	outStream io.Writer
	errStream io.Writer
}

func (c *CLI) Run(args []string) int {
	var userID string

	flags := flag.NewFlagSet("gorden", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.StringVar(&userID, "userID", getEnvString(envUserID, ""), "github account")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagError
	}

	if args[1] != "-userID" {
		userID = args[1]
	}

	if userID == "" {
		flags.Usage()
		return ExitCodeParseFlagError
	}

	url := "https://github.com/users/" + userID + "/contributions"

	contributions, sum := GetContributions(url)

	fmt.Fprintf(c.outStream, "Current Streak: %d\n", contributions)
	fmt.Fprintf(c.outStream, "Year of contributions: %d\n", sum)

	return ExitCodeOK
}

func GetContributions(url string) (int, int) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(err)
	}

	var contributions, sum int

	doc.Find("rect").Each(func(_ int, s *goquery.Selection) {
		datacount, _ := s.Attr("data-count")
		c, _ := strconv.Atoi(datacount)

		if c > 0 {
			contributions++
		} else {
			contributions = 0
		}

		sum += c
	})

	return contributions, sum
}

func getEnvString(env, def string) string {
	r := os.Getenv(env)

	if r == "" {
		return def
	}

	return r
}
