// Package image provides ...
package image

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	// "golang.org/x/image/webp"
)

var ErrUnsupported = errors.New("unsupported image type")

func Encode(fm Format, w io.Writer, m image.Image) (err error) {
	switch fm {
	case GIF:
		err = gif.Encode(w, m, nil)
	case JPG, JPEG:
		err = jpeg.Encode(w, m, nil)
	case PNG:
		err = png.Encode(w, m)
	case TIFF:
		err = tiff.Encode(w, m, nil)
	case BMP:
		err = bmp.Encode(w, m)
	default:
		return ErrUnsupported
	}

	return
}
