package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var Usage string = `
Keep downloads and organizes repositories.

Usage:

	keep [arguments] <repository>
`

type githubContent struct {
	organization string
}

type githubResponse struct {
	FullName string `json:"full_name"`
}

// die prints error then exit.
func die(err interface{}) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	var fork bool
	flag.BoolVar(&fork, "fork", false, "create a fork of the repo, then download the fork")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, Usage)
		flag.PrintDefaults()
		os.Exit(1)
	}
	addr := args[0]
	if strings.Contains(addr, "://") {
		die("keep always use https so you don't need to specify")
	}
	hostPath := strings.SplitN(addr, "/", 2)
	if len(hostPath) != 2 {
		die("invalid form of address:" + addr)
	}
	host := hostPath[0]
	path := hostPath[1]
	switch host {
	case "github.com":
		user := os.Getenv("KEEP_GITHUB_USER")
		token := os.Getenv("KEEP_GITHUB_AUTH")
		if fork {
			if user == "" {
				die("KEEP_GITHUB_USER should be specified to fork")
			}
			if token == "" {
				die("KEEP_GITHUB_AUTH should be specified to fork")
			}
		}

		paths := strings.Split(path, "/")
		if len(paths) != 2 {
			die("invalid repository path:" + path)
		}
		org := paths[0]
		repo := paths[1]

		upstream := ""
		origin := addr
		if fork {
			upstream = addr
			// origin will set by fork api response
		}

		if fork {
			// request a fork of the repo
			forkApiAddr := fmt.Sprintf("https://api.github.com/repos/%s/%s/forks", org, repo)
			content, err := json.Marshal(githubContent{organization: user})
			if err != nil {
				die(fmt.Errorf("unable to marshal githubContent: %w", err))
			}
			req, err := http.NewRequest("POST", forkApiAddr, bytes.NewBuffer(content))
			if err != nil {
				die(err)
			}
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Accept", "application/vnd.github.v3+json")
			req.Header.Add("Authorization", "token "+token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				die(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				die(fmt.Errorf("could not read response body: %w", err))
			}
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				die(fmt.Sprintf("bad response status: %d\n%s", resp.StatusCode, string(body)))
			}
			// successfully forked, or it has existed already.
			forked := githubResponse{}
			json.Unmarshal(body, &forked)
			origin = "github.com/" + forked.FullName
		}

		keepPath := os.Getenv("KEEPPATH")
		if keepPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				die(err)
			}
			keepPath = home + "/src"
		}
		dst := keepPath + "/" + origin
		_, err := os.Stat(dst)
		if err != nil && !os.IsNotExist(err) {
			die(err)
		} else if err == nil {
			die(fmt.Sprintf("dest directory already exists: %s", dst))
		}
		dstParent := filepath.Dir(dst)
		err = os.MkdirAll(dstParent, 0755)
		if err != nil {
			die(fmt.Errorf("could not make intermediate directories: %w", err))
		}
		cloneAddr := origin
		if token != "" {
			cloneAddr = token + "@" + cloneAddr
		} else if user != "" {
			cloneAddr = user + "@" + cloneAddr
		}
		cmd := exec.Command("git", "clone", "https://"+cloneAddr, dst)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			die(err)
		}
		if upstream != "" {
			if token != "" {
				upstream = token + "@" + upstream
			} else if user != "" {
				upstream = user + "@" + upstream
			}
			cmd = exec.Command("git", "remote", "add", "upstream", "https://"+upstream)
			cmd.Dir = dst
			out, err := cmd.CombinedOutput()
			if err != nil {
				die(fmt.Sprintf("%s\n%s", string(out), err))
			}
		}
	default:
		die(fmt.Sprintf("unsupported host: %s", host))
	}

}
