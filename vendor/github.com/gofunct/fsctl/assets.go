package fsctl

import (
	"github.com/Masterminds/sprig"
	"io"
	"path/filepath"
	"text/template"
)

type AssetFunc func(string) ([]byte, error)
type AssetDirFunc func(string) ([]string, error)

func (f *Fs) SetAssetProcessor(a AssetFunc) {
	f.assetFunc = a
}

func (f *Fs) SetAssetDirProcessor(a AssetDirFunc) {
	f.dirFunc = a
}

func (f *Fs) loadDirectory(directory string) (*template.Template, error) {
	var tmpl *template.Template

	return f.assetToTmpl(tmpl, directory)
}

func (f *Fs) loadHtmlDirectory(directory string) (*template.Template, error) {
	var tmpl *template.Template
	return f.htmlassetToTmpl(tmpl, directory)
}

func (f *Fs) assetToTmpl(tmpl *template.Template, directory string) (*template.Template, error) {
	files, err := f.dirFunc(directory)
	if err != nil {
		return tmpl, err
	}

	for _, filePath := range files {
		contents, err := f.assetFunc(directory + "/" + filePath)
		if err != nil {
			return tmpl, err
		}

		name := filepath.Base(filePath)

		if tmpl == nil {
			tmpl = template.New(name)
		}

		if name != tmpl.Name() {
			tmpl = tmpl.New(name)
		}
		tmpl.Funcs(sprig.GenericFuncMap())

		if _, err = tmpl.Parse(string(contents)); err != nil {
			return tmpl, err
		}
	}

	return tmpl, nil
}

func (f *Fs) htmlassetToTmpl(tmpl *template.Template, directory string) (*template.Template, error) {
	files, err := f.dirFunc(directory)
	if err != nil {
		return tmpl, err
	}

	for _, filePath := range files {
		contents, err := f.assetFunc(directory + "/" + filePath)
		if err != nil {
			return tmpl, err
		}

		name := filepath.Base(filePath)

		if tmpl == nil {
			tmpl = template.New(name)
		}

		if name != tmpl.Name() {
			tmpl = tmpl.New(name)
		}
		tmpl.Funcs(sprig.GenericFuncMap())

		if _, err = tmpl.Parse(string(contents)); err != nil {
			return tmpl, err
		}
	}

	return tmpl, nil
}

func (f *Fs) MustParseAssets(directory string) *template.Template {
	if tmpl, err := f.loadDirectory(directory); err != nil {
		panic(err)
	} else {
		return tmpl
	}
}

func (f *Fs) MustParseHtmlAssets(directory string) *template.Template {
	if tmpl, err := f.loadHtmlDirectory(directory); err != nil {
		panic(err)
	} else {
		return tmpl
	}
}

// Template reads a go template and writes it to dist given data.
func (f *Fs) MustExecAssets(dir string, w io.Writer) error {
	tmpl := f.MustParseAssets(dir)
	return tmpl.Execute(w, f.AllSettings())
}

// Template reads a go template and writes it to dist given data.
func (f *Fs) MustExecHtmlAssets(dir string, w io.Writer) error {
	tmpl := f.MustParseHtmlAssets(dir)
	return tmpl.Execute(w, f.AllSettings())
}
