package fs

import (
	"bufio"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/jessevdk/go-assets"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Fs struct {
	*afero.Afero
	*viper.Viper
}

func NewFs() *Fs {
	fs := &afero.Afero{
		Fs: afero.NewOsFs(),
	}
	v := viper.GetViper()
	if v == nil {
		v = viper.New()
	}
	v.AutomaticEnv()
	v.SetFs(fs)
	
	f :=  &Fs{
		Afero:      fs,
		Viper:      v,
	}
	f.Sync()
	return f
}

// Template reads a go template and writes it to dist given data.
func (f *Fs) ProcessAsset(t *template.Template, file *assets.File) {
	if file.Name() == "/" {
		return
	}
	content := string(file.Data)

	tpl := t.New(file.Name()).Funcs(sprig.GenericFuncMap())
	tpl, err := tpl.Parse(string(content))
	if err != nil {
		f.Panic(err, "Could not parse template ")
	}

	fl, err := f.Create(file.Name())
	if err != nil {
		f.Panic(err, "Could not create file for writing")
	}
	defer fl.Close()
	err = tpl.Execute(fl, f.AllSettings())
	if err != nil {
		f.Panic(err, "Could not execute template")
	}
}

func (f *Fs)  WalkTemplates(dir string, outDir string) {

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			f.Panic(err, "error walking path")
		}
		if strings.Contains(path, ".tmpl") {
			b, err := ioutil.ReadFile(path)
			newt, err := template.New(info.Name()).Funcs(sprig.GenericFuncMap()).Parse(string(b))
			if err != nil {
				return err
			}

			f, err := f.Afero.Create(outDir + "/" + strings.TrimSuffix(info.Name(), ".tmpl"))
			if err != nil {
				return err
			}
			return newt.Execute(f, f.v.AllSettings())
		}
		return nil
	}); err != nil {
		f.Panic(err, "failed to walk templates")
	}
}

func (f *Fs)  CopyFile(srcfile, dstfile string) (*afero.File, error) {
	srcF, err := f.Open(srcfile) // nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("could not open source file: %s", err)
	}
	defer srcF.Close()

	dstF, err := f.Afero.Create(dstfile)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(dstF, srcF); err != nil {
		return nil, fmt.Errorf("could not copy file: %s", err)
	}
	return &dstF, f.Chmod(dstfile, 0755)
}

func (c *Fs)  ScanAndReplaceFile(f afero.File, replacements ...string) {
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

func  (f *Fs)  ScanAndReplace(r io.Reader, replacements ...string) string {
	scanner := bufio.NewScanner(r)
	rep := strings.NewReplacer(replacements...)
	var text string
	for scanner.Scan() {
		text = rep.Replace(scanner.Text())
	}
	return text
}
