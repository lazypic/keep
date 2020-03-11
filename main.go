package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// die prints error then exit.
func die(err interface{}) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		die("clone address not specified")
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
		paths := strings.Split(path, "/")
		if len(paths) != 2 {
			die("invalid repository path:" + path)
		}

		keepPath := os.Getenv("KEEPPATH")
		if keepPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				die(err)
			}
			keepPath = home + "/src"
		}
		dst := keepPath + "/" + addr
		_, err := os.Stat(dst)
		if err != nil {
			if !os.IsNotExist(err) {
				die(err)
			}
		} else if err == nil {
			die(fmt.Sprintf("dest directory already exists: %s", dst))
		}
		err = os.MkdirAll(dst, 0755)
		if err != nil {
			die(fmt.Errorf("could not make intermediate directories: %w", err))
		}
		cmd := exec.Command("git", "clone", "https://"+addr, dst)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			die(err)
		}

		// you don't need permission for 'origin' except some specific times.
		// when you need it, use remote 'origin-private' instead.
		// which is actually 'origin' with the authentication token.
		token := os.Getenv("KEEP_GITHUB_AUTH")
		if token == "" {
			return
		}
		cmd = exec.Command("git", "-C", dst, "remote", "add", "origin-private", "https://"+token+"@"+addr)
		b, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprint(os.Stderr, string(b))
			die(err)
		}
	default:
		die(fmt.Sprintf("unsupported host: %s", host))
	}
}
