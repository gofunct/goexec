package goexec

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/jessevdk/go-assets"
	"github.com/spf13/afero"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Template reads a go template and writes it to dist given data.
func (c *Command) ProcessAsset(t *template.Template, file *assets.File) {
	if file.Name() == "/" {
		return
	}
	content := string(file.Data)

	tpl := t.New(file.Name()).Funcs(sprig.GenericFuncMap())
	tpl, err := tpl.Parse(string(content))
	if err != nil {
		c.Panic(err, "Could not parse template ")
	}

	f, err := c.Create(file.Name())
	if err != nil {
		c.Panic(err, "Could not create file for writing")
	}
	defer f.Close()
	err = tpl.Execute(f, c.v.AllSettings())
	if err != nil {
		c.Panic(err, "Could not execute template")
	}
}

func (c *Command) WalkTemplates(dir string, outDir string) {

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			c.Panic(err, "error walking path")
		}
		if strings.Contains(path, ".tmpl") {
			b, err := ioutil.ReadFile(path)
			newt, err := template.New(info.Name()).Funcs(sprig.GenericFuncMap()).Parse(string(b))
			if err != nil {
				return err
			}

			f, err := c.Afero.Create(outDir + "/" + strings.TrimSuffix(info.Name(), ".tmpl"))
			if err != nil {
				return err
			}
			return newt.Execute(f, c.v.AllSettings())
		}
		return nil
	}); err != nil {
		c.Panic(err, "failed to walk templates")
	}
}

func (c *Command) CopyFile(srcfile, dstfile string) (*afero.File, error) {
	srcF, err := c.Open(srcfile) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("could not open source file: %s", err)
	}
	defer srcF.Close()

	dstF, err := c.Afero.Create(dstfile)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(dstF, srcF); err != nil {
		return nil, fmt.Errorf("could not copy file: %s", err)
	}
	return &dstF, c.Chmod(dstfile, 0755)
}

func (c *Command) ScanAndReplaceFile(f afero.File, replacements ...string) {
	nm := f.Name()
	d, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err.Error())
	}
	if err := c.Remove(f.Name()); err != nil {
		panic(err.Error())
	}
	scanner := bufio.NewScanner(strings.NewReader(fmt.Sprintf("%s", d)))
	rep := strings.NewReplacer(replacements...)
	var newstr string
	for scanner.Scan() {
		newstr = rep.Replace(scanner.Text())
		if err := scanner.Err(); err != nil {
			fmt.Println(err.Error())
			break
		}
	}
	newf, err := c.Create(nm)
	if err != nil {
		panic(err.Error())
	}
	_, err = io.WriteString(newf, newstr)
	c.Panic(err, "failed to write string to new file")
	c.Println("successfully scanned and replaced: " + f.Name())

}

func (c *Command) ScanAndReplace(r io.Reader, replacements ...string) string {
	scanner := bufio.NewScanner(r)
	rep := strings.NewReplacer(replacements...)
	var text string
	for scanner.Scan() {
		text = rep.Replace(scanner.Text())
	}
	return text
}
