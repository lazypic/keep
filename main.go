package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type githubContent struct {
	organization string
}

func main() {
	args := os.Args[1:]
	addr := args[0]
	if strings.Contains(addr, "://") {
		fmt.Fprintln(os.Stderr, "fork always use https so you don't need to specify:")
		os.Exit(1)
	}
	hostPath := strings.SplitN(addr, "/", 2)
	if len(hostPath) != 2 {
		fmt.Fprintln(os.Stderr, "invalid form of address")
		os.Exit(1)
	}
	host := hostPath[0]
	path := hostPath[1]
	switch host {
	case "github.com":
		user := os.Getenv("FORK_GITHUB_USER")
		if user == "" {
			fmt.Fprintln(os.Stderr, "FORK_GITHUB_USER not defined")
			os.Exit(1)
		}
		token := os.Getenv("FORK_GITHUB_AUTH")
		if token == "" {
			fmt.Fprintln(os.Stderr, "FORK_GITHUB_AUTH not defined")
			os.Exit(1)
		}
		paths := strings.Split(path, "/")
		if len(paths) != 2 {
			fmt.Fprintln(os.Stderr, "invalid repository address")
			os.Exit(1)
		}
		org := paths[0]
		repo := paths[1]
		forkApiAddr := fmt.Sprintf("https://api.github.com/repos/%s/%s/forks", org, repo)
		content, err := json.Marshal(githubContent{organization: user})
		if err != nil {
			fmt.Fprintln(os.Stderr, "unable to marshal githubContent")
		}
		req, err := http.NewRequest("POST", forkApiAddr, bytes.NewBuffer(content))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "token "+token)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		_ = body
		// fmt.Println(string(body))
	default:
		fmt.Fprintln(os.Stderr, "unsupported host:", host)
		os.Exit(1)
	}

}
