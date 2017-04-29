// Package mux provides ...
package mux

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	// 20 MB
	defaultMaxMemory = 20 << 20
)

// simple context
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
}

// reply with string
func (c *Context) String(code int, val string) {
	c.Status(code)
	c.Writer.Write([]byte(val))
}

// reply with json
func (c *Context) Json(code int, val interface{}) error {
	c.Status(code)
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(c.Writer).Encode(val)
}

// reply with binary stream
func (c *Context) Stream(code int, val []byte) {
	c.Status(code)
}

// reply with file by download
func (c *Context) File(code int, name string, val []byte) {
	c.Status(code)

	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=\""+name+"\"")
	c.Writer.Write(val)
}

// get upload file by key
func (c *Context) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Request.FormFile(key)
}

// set http code
func (c *Context) Status(code int) {
	if code > 0 {
		c.Writer.WriteHeader(code)
	}
}

// get upload files
func (c *Context) FormFiles() ([]*multipart.FileHeader, error) {
	// max size 20 MB
	err := c.Request.ParseMultipartForm(defaultMaxMemory)
	if err != nil {
		return nil, err
	}

	if c.Request.MultipartForm == nil || c.Request.MultipartForm.File == nil {
		return nil, http.ErrMissingFile
	}

	var fhs []*multipart.FileHeader
	for _, fs := range c.Request.MultipartForm.File {
		fhs = append(fhs, fs...)
	}

	return fhs, nil
}

// get param
func (c *Context) Query(key string) string {
	if vals, ok := c.Request.URL.Query()[key]; ok && len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// post param
func (c *Context) FormValue(key string) string {
	return c.Request.FormValue(key)
}

// string param with default
func (c *Context) DefaultString(key, d string) string {
	v := c.Request.FormValue(key)
	if v == "" {
		return d
	}
	return v
}

// int param
func (c *Context) GetInt(key string) (int, error) {
	val := c.Request.FormValue(key)
	return strconv.Atoi(val)
}

// int param with default
func (c *Context) DefaultInt(key string, d int) int {
	v, err := c.GetInt(key)
	if err != nil {
		return d
	}
	return v
}

// bool param with default
func (c *Context) DefaultBool(key string, b bool) bool {
	val := c.Request.FormValue(key)
	v, err := strconv.ParseBool(val)
	if err != nil {
		return b
	}
	return v
}

// owner handlefunc
type HandleFunc func(*Context)

// implement handler
func (h HandleFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{
		Writer:  w,
		Request: r,
	}

	h(c)
}
