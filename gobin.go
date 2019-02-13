package goexec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/rogpeppe/go-internal/module"
)

const (
	tempModule = "temporary.com/gobin"
)

func (c *Command) GoBin(pkgPatts []string) error {

	var gopath string          // effective GOPATH
	var gobinCache string      // does what it says on the tin
	var localCacheProxy string // local filesystem-based module download cache

	// cache path discovery
	{
		gopath = os.Getenv("GOPATH")
		if gopath != "" {
			gopath = filepath.SplitList(gopath)[0]
		} else {
			uhd, err := userHomeDir()
			if err != nil {
				return fmt.Errorf("failed to determine user home directory: %v", err)
			}
			gopath = filepath.Join(uhd, "go")
		}

		// TODO I don't think the module cache path is advertised anywhere public...
		// intentionally but in case it is, replace what follows
		localCacheProxy = "GOPROXY=file://" + path.Join(filepath.ToSlash(gopath), "pkg", "mod", "cache", "download")

		ucd, err := os.UserCacheDir()
		if err != nil {
			return fmt.Errorf("failed to determine user cache dir: %v", err)
		}

		gobinCache = filepath.Join(ucd, "gobin")
	}

	var allPkgs []*arg // all of the non-run command line provided packages
	var netPkgs []*arg // packages that need network resolution
	var nonMain []*arg // non-main packages

	// prepare allPkgs
	{
		if len(pkgPatts) == 0 {
			return fmt.Errorf("need to provide at least one main package")
		}

		var tmpDirs []string
		defer func() {
			for _, td := range tmpDirs {
				os.RemoveAll(td)
			}
		}()

		for _, patt := range pkgPatts {
			parts := strings.SplitN(patt, "@", 2)

			a := &arg{
				patt:    patt,
				pkgPatt: parts[0],
			}

			if len(parts) == 2 {
				a.verPatt = parts[1]
			}

			td, err := ioutil.TempDir("", "gobin")
			if err != nil {
				return fmt.Errorf("failed to create temp dir: %v", err)
			}
			tmpDirs = append(tmpDirs, td)
			if err := ioutil.WriteFile(filepath.Join(td, "go.mod"), []byte("module "+tempModule+"\n"), 0644); err != nil {
				return fmt.Errorf("failed to initialise temp Go module: %v", err)
			}
			a.wd = td

			allPkgs = append(allPkgs, a)
		}
	}
	netPkgs = allPkgs

	// network resolution step
	for _, pkg := range netPkgs {
		proxy := os.Getenv("GOPROXY")
		if proxy != "" {
			proxy = "GOPROXY=" + proxy
		}

		if err := pkg.get(proxy); err != nil {
			return err
		}

		if err := pkg.list(proxy); err != nil {
			return err
		}

		if pkg.resErr != nil {
			nonMain = append(nonMain, pkg)
		}
	}

	if len(nonMain) > 0 {
		for _, pkg := range nonMain {
			fmt.Fprintf(os.Stderr, "%v@%v: %v\n", pkg.pkgPatt, pkg.verPatt, pkg.resErr)
		}
		s := ""
		if len(nonMain) > 1 {
			s = "s"
		}
		return fmt.Errorf("failed to resolve module-based main package%v", s)
	}

	for _, pkg := range allPkgs {
		// each mainPkg install must be done as a separate go command invocation because
		// we set a different GOBIN for each one.
		for _, mp := range pkg.mainPkgs {
			// calculate the relative install directory from main package import path
			// and the containing module's version
			var mainrel string
			{
				emp, err := module.EncodePath(mp.Module.Path)
				if err != nil {
					return fmt.Errorf("failed to encode module path %v: %v", mp.Module.Path, err)
				}

				md := filepath.FromSlash(emp)

				if mp.Module.Version != "" {
					emv, err := module.EncodeVersion(mp.Module.Version)
					if err != nil {
						return fmt.Errorf("failed to encode module version %v: %v", mp.Module.Version, err)
					}

					md = filepath.Join(md, "@v", emv)
				}

				epp, err := module.EncodePath(mp.ImportPath)
				if err != nil {
					return fmt.Errorf("failed to encode package relative path %v: %v", mp.ImportPath, err)
				}
				mainrel = filepath.Join(md, filepath.FromSlash(epp))
			}

			gobin := filepath.Join(gobinCache, mainrel)
			target := filepath.Join(gobin, path.Base(mp.ImportPath))

			if runtime.GOOS == "windows" {
				target += ".exe"
			}

			// optimistically remove our target in case we are installing over self
			// TODO work out what to do for Windows
			if mp.ImportPath == "github.com/myitcv/gobin" {
				_ = os.Remove(target)
			}

			var stdout bytes.Buffer

			installCmd := exec.Command("go", "install", mp.ImportPath)
			installCmd.Dir = pkg.wd
			installCmd.Env = append(buildEnv(localCacheProxy), "GOBIN="+gobin)
			installCmd.Stdout = &stdout

			if err := run(installCmd); err != nil {
				return err
			}
			fmt.Println(target)
			fmt.Printf("%v %v\n", mp.Module.Path, mp.Module.Version)
			installBin := os.Getenv("GOBIN")
			if installBin == "" {
				installBin = filepath.Join(gopath, "bin")
			}
			if err := os.MkdirAll(installBin, 0755); err != nil {
				return fmt.Errorf("failed to mkdir %v: %v", installBin, err)
			}
			src, err := os.Open(target)
			if err != nil {
				return fmt.Errorf("failed to open %v: %v", target, err)
			}
			defer src.Close()
			bin := filepath.Join(installBin, path.Base(mp.ImportPath))

			if runtime.GOOS == "windows" {
				bin += ".exe"
			}

			openMode := os.O_CREATE | os.O_WRONLY

			// optimistically remove our target in case we are installing over self
			// TODO work out what to do for Windows
			if mp.ImportPath == "github.com/myitcv/gobin" {
				_ = os.Remove(bin)
				openMode = openMode | os.O_EXCL
			}

			dest, err := os.OpenFile(bin, openMode, 0755)
			if err != nil {
				return fmt.Errorf("failed to open %v for writing: %v", bin, err)
			}
			defer dest.Close()
			if _, err := io.Copy(dest, src); err != nil {
				return fmt.Errorf("failed to copy %v to %v", target, bin)
			}
			fmt.Printf("Installed %v@%v to %v\n", mp.ImportPath, mp.Module.Version, bin)
		}
	}

	return nil
}

// listPkg is a convenience type for unmarshaling the output from go list
type listPkg struct {
	ImportPath string
	Name       string
	Dir        string
	Module     struct {
		Path    string
		Dir     string
		Version string
		GoMod   string
	}
}

// arg is a wrapper around a command line-provided package
type arg struct {
	patt     string     // the command line-provided pattern
	pkgPatt  string     // the package part of patt
	verPatt  string     // the version part of patt
	mainPkgs []*listPkg // main packages resolved from patt
	wd       string     // working directory for resolution
	resErr   error      // resolution error
	target   string     // the gobin cache target
}

var (
	errNonMain      = errors.New("not a main package")
	errMultiModules = errors.New("cannot (yet) install main packages from a pattern that spans multiple modules")
)

// resolve attempts to resolve a.patt to main packages, using the supplied
// proxy (if != "").  If there is an error resolving a.patt to a package and
// version this is returned. Otherwise the main packages matched by the
// packages are populated into a.mainPkgs
func (a *arg) get(proxy string) error {
	env := buildEnv(proxy)

	getCmd := exec.Command("go", "get", "-d", a.patt)
	getCmd.Dir = a.wd
	getCmd.Env = env

	if err := run(getCmd); err != nil {
		return err
	}

	return nil
}

func (a *arg) list(proxy string) error {
	env := buildEnv(proxy)

	var stdout bytes.Buffer

	listCmd := exec.Command("go", "list", "-json", a.pkgPatt)
	listCmd.Dir = a.wd
	listCmd.Stdout = &stdout
	listCmd.Env = env

	if err := run(listCmd); err != nil {
		return err
	}

	dec := json.NewDecoder(&stdout)

	// TODO if/when we support patterns including ... we will need to change the
	// semantics of a.resErr and the version resolution below

	// TODO for now we simply throw an error in case the package pattern
	// provided cross module boundaries. Because in global mode, we would need a
	// temp module per module for things to work cleanly (else we might have to
	// handle conflicting replace statements)
	seenMods := make(map[string]bool)

	for {
		pkg := new(listPkg)
		if err := dec.Decode(pkg); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		a.verPatt = pkg.Module.Version

		if pkg.Name != "main" {
			a.resErr = errNonMain
			return nil
		}

		seen := seenMods[pkg.Module.Path]
		if !seen && len(seenMods) > 0 {
			a.resErr = errMultiModules
			return nil
		}

		seenMods[pkg.Module.Path] = true

		// If we are not in main-module mode (i.e. -m is not provided), then we
		// are working in a temporary module. Any replacements in
		// $(mp.Module.GoMod) need to be applied to $(pkg.wd)/go.mod. So the
		// simplest thing to do is copy the main package's module's go.mod over
		// the top of the temporary module's go.mod and then adjust the module
		// line and add the single requirement that we now have resolved.
		if !seen {
			srcPath := pkg.Module.GoMod
			src, err := os.Open(srcPath)
			if err != nil {
				return fmt.Errorf("failed to open src %v: %v", srcPath, err)
			}

			destPath := filepath.Join(a.wd, "go.mod")
			dest, err := os.Create(destPath)
			if err != nil {
				return fmt.Errorf("failed to create dest %v: %v", destPath, err)
			}

			_, err = io.Copy(dest, src)
			src.Close()
			if err != nil {
				return fmt.Errorf("failed to copy %v to %v: %v", srcPath, destPath, err)
			}
			if err := dest.Close(); err != nil {
				return fmt.Errorf("failed to close dest file %v: %v", destPath, err)
			}

			// work around https://github.com/golang/go/issues/28820 by reading the go.mod
			// and doing a string replace on the module line. Otherwise we could do this
			// in the go mod edit below.
			{
				modreg := regexp.MustCompile(`\A\s*module\s+"?` + regexp.QuoteMeta(pkg.Module.Path) + `"?.*\n`)
				fpath := filepath.Join(a.wd, "go.mod")
				fbyts, err := ioutil.ReadFile(fpath)
				if err != nil {
					return fmt.Errorf("failed to read %v: %v", fpath, err)
				}
				fstr := string(fbyts)
				fstr = modreg.ReplaceAllString(fstr, "module "+tempModule+"\n")
				if err := ioutil.WriteFile(fpath, []byte(fstr), 0644); err != nil {
					return fmt.Errorf("failed to write back to %v: %v", fpath, err)
				}
			}

			gmeCmd := exec.Command("go", "mod", "edit", "-require="+pkg.Module.Path+"@"+pkg.Module.Version)
			gmeCmd.Dir = a.wd
			gmeCmd.Env = buildEnv("")

			if err := run(gmeCmd); err != nil {
				return err
			}

			// now that we effectively have a copy of everything relevant in the
			// target module (including replace directives), list to ensure they
			// have been resolved

			listCmd := exec.Command("go", "list", "-json", pkg.ImportPath)
			listCmd.Dir = a.wd
			listCmd.Env = buildEnv(proxy)

			if err := run(listCmd); err != nil {
				return err
			}
		}

		a.mainPkgs = append(a.mainPkgs, pkg)
	}

	return nil
}

// buildEnv builds the correct environment for running go commands from gobin.
// proxy is expected to be empty or take the form "GOPROXY=X". If it is non
// empty it will be added to the environment.
func buildEnv(proxy string) []string {
	env := append(os.Environ(), "GO111MODULE=on")
	if proxy != "" {
		env = append(env, proxy)
	}
	goflags := os.Getenv("GOFLAGS")
	return append(env, "GOFLAGS="+goflags)
}

func run(cmd *exec.Cmd) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	start := time.Now()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run %v: %v\n%s", strings.Join(cmd.Args, " "), err, stderr.String())
	}

	end := time.Now()

	env := cmd.Env
	if env == nil {
		env = os.Environ()
	}
	var goenv []string
	for _, v := range env {
		if strings.HasPrefix(v, "GO") {
			goenv = append(goenv, v)
		}
	}
	fmt.Fprintf(os.Stderr, "+ cd %v; %v %v # took %v\n%s", cmd.Dir, strings.Join(goenv, " "), strings.Join(cmd.Args, " "), end.Sub(start), stderr.String())

	return nil
}

// os.UserHomeDir() is implemented in
// https://go-review.googlesource.com/c/go/+/139418 so implement here for now
func userHomeDir() (string, error) {
	env, enverr := "HOME", "$HOME"
	switch runtime.GOOS {
	case "windows":
		env, enverr = "USERPROFILE", "%userprofile%"
	case "plan9":
		env, enverr = "home", "$home"
	case "nacl", "android":
		return "/", nil
	case "darwin":
		if runtime.GOARCH == "arm" || runtime.GOARCH == "arm64" {
			return "/", nil
		}
	}
	if v := os.Getenv(env); v != "" {
		return v, nil
	}
	return "", errors.New(enverr + " is not defined")
}
