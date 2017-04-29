// Package services provides ...
package services

import (
	"archive/zip"
	"net/http"

	"github.com/imgee/imgee/httpd/mux"
	"github.com/imgee/imgee/resize"
)

// resize api handler
func HandleResize(c *mux.Context) {
	// get params
	itpl := c.DefaultString("itpl", "l3")
	w := c.DefaultInt("w", 0)
	h := c.DefaultInt("h", 0)
	thm := c.DefaultBool("thm", false)

	// zip file
	zipper := zip.NewWriter(c.Writer)
	defer zipper.Close()

	rs := &resize.Resize{
		Width:       w,
		Height:      h,
		Interp:      itpl,
		IsThumbnail: thm,
	}

	// get all upload files
	fhs, err := c.FormFiles()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// resize all image
	for _, fh := range fhs {
		f, err := fh.Open()
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		defer f.Close()

		zf, err := zipper.Create(fh.Filename)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		err = rs.Resize(f, zf)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
	}

	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=\"imgee.zip\"")
}
