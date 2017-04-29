// Package resize provides ...
package resize

import (
	"archive/zip"
	"errors"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"

	img "github.com/imgee/image"
	ext "github.com/imgee/image/extension"
	"github.com/imgee/imgee/utils"
	rsz "github.com/imgee/resize"
)

const (
	command = "resize"
	usage   = `usage:
	resize [thumbnail] path [options]
options:
	-i interp	Interpolation function in "nn", "bc", "bl", "mn", "l2", "l3" (default: l3)
	-w pixel	The width of resized image (default 0)
	-h pixel	The height of resized image (default 0)
	-o path		The resized image saved as...
Example:
	resize -o static/image/a.zip test.png
	resize -i l2 -w 200 test.png
	resize -i l3 -w 200 -h 100 -o static/image/ test.png`
)

// Interpolation Function
var Interpolation = map[string]rsz.InterpolationFunction{
	"nn": rsz.NearestNeighbor,
	"bc": rsz.Bicubic,
	"bl": rsz.Bilinear,
	"mn": rsz.MitchellNetravali,
	"l2": rsz.Lanczos2,
	"l3": rsz.Lanczos3,
}

// some configuration before resize
type Resize struct {
	Out         io.Writer
	Paths       []string
	IsThumbnail bool
	Width       int
	Height      int
	Interp      string
}

// terminal exec function
func (rs *Resize) Exec(args string) error {
	// parse params
	nArgs, flag := utils.Flag(args)
	// print help
	if _, ok := flag["help"]; ok ||
		len(nArgs) == 0 ||
		nArgs == "thumbnail" {
		rs.help()
		return nil
	}

	// second command and path
	fields := strings.Fields(nArgs)
	var paths []string
	for _, v := range fields {
		if v == "thumbnail" {
			rs.IsThumbnail = true
			continue
		}
		paths = append(paths, v)
	}
	rs.Paths = paths

	// flags
	rs.Interp = utils.SetFlag(flag, "i", "l3").(string)
	rs.Width = utils.SetFlag(flag, "w", 0).(int)
	rs.Height = utils.SetFlag(flag, "h", 0).(int)

	// output path
	out := utils.SetFlag(flag, "o", "imgee.zip").(string)
	dir, file := filepath.Split(out)
	if file == "" {
		file = "imgee.zip"
	}

	f, err := os.OpenFile(filepath.Join(dir, file), os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	rs.Out = f

	// handle with image
	return rs.walk()
}

// sub commands
func (rs *Resize) Cmds() []string {
	return []string{
		"thumbnail",
	}
}

// command
func (rs *Resize) Command() string {
	return command
}

// print usage
func (rs *Resize) help() {
	print(usage)
}

// walk files with given dir
func (rs *Resize) walk() error {
	zipper := zip.NewWriter(rs.Out)
	defer zipper.Close()

	for _, p := range rs.Paths {
		// create zip file
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		// file info
		info, err := f.Stat()
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		w, err := zipper.CreateHeader(header)
		if err != nil {
			return err
		}

		// resize
		err = rs.Resize(f, w)
		if err != nil {
			return err
		}

		// write file
		_, err = io.Copy(w, f)
		if err != nil {
			return err
		}
	}

	return nil
}

// resize image by config
func (rs *Resize) Resize(r io.Reader, w io.Writer) error {
	m, name, err := image.Decode(r)
	if err != nil {
		return nil
	}

	fm, ok := ext.Extensions[name]
	if !ok {
		return errors.New("unknown format")
	}

	interp, ok := Interpolation[rs.Interp]
	if !ok {
		return errors.New("unknown interpolation")
	}

	if rs.IsThumbnail {
		m = rsz.Thumbnail(uint(rs.Width), uint(rs.Height), m, interp)
		if err != nil {
			return err
		}
	} else {
		m = rsz.Resize(uint(rs.Width), uint(rs.Height), m, interp)
		if err != nil {
			return err
		}
	}

	return img.Encode(fm, w, m)
}
