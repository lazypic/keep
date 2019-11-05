package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type githubContent struct {
	organization string
}

// die prints msg then exit.
func die(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	args := os.Args[1:]
	addr := args[0]
	if strings.Contains(addr, "://") {
		die("fork always use https so you don't need to specify")
	}
	hostPath := strings.SplitN(addr, "/", 2)
	if len(hostPath) != 2 {
		die("invalid form of address")
	}
	host := hostPath[0]
	path := hostPath[1]
	switch host {
	case "github.com":
		user := os.Getenv("FORK_GITHUB_USER")
		if user == "" {
			die("FORK_GITHUB_USER not defined")
		}
		token := os.Getenv("FORK_GITHUB_AUTH")
		if token == "" {
			die("FORK_GITHUB_AUTH not defined")
		}
		paths := strings.Split(path, "/")
		if len(paths) != 2 {
			die("invalid repository address")
		}
		org := paths[0]
		repo := paths[1]
		forkApiAddr := fmt.Sprintf("https://api.github.com/repos/%s/%s/forks", org, repo)
		content, err := json.Marshal(githubContent{organization: user})
		if err != nil {
			die("unable to marshal githubContent")
		}
		req, err := http.NewRequest("POST", forkApiAddr, bytes.NewBuffer(content))
		if err != nil {
			die(err.Error())
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/vnd.github.v3+json")
		req.Header.Add("Authorization", "token "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			die(err.Error())
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			die(err.Error())
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			die(fmt.Sprintf("bad reponse status: %s\n%s", resp.StatusCode, string(body)))
		}

		// successfully forked, or it has existed already.
		home, err := os.UserHomeDir()
		if err != nil {
			die(err.Error())
		}
		dst := home + "/src/" + host + "/" + user + "/" + repo
		if err != nil {
			die(err.Error())
		}
		_, err = os.Stat(dst)
		if err != nil && !os.IsNotExist(err) {
			die(err.Error())
		} else if err == nil {
			die(fmt.Sprintf("dest directory already exist: %s", dst))
		}
		dstParent := filepath.Dir(dst)
		err = os.MkdirAll(dstParent, 0755)
		cmd := exec.Command("git", "clone", "https://"+addr, dst)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			die(err.Error())
		}
		cmd = exec.Command("git", "remote", "add", "upstream", "https://"+addr)
		cmd.Dir = dst
		out, err := cmd.CombinedOutput()
		if err != nil {
			die(fmt.Sprintf("%s\n%s", string(out), err))
		}
	default:
		die(fmt.Sprintf("unsupported host: %s", host))
	}

}
