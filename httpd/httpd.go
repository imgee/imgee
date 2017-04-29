// Package httpd provides ...
package httpd

import (
	"fmt"
	"net/http"

	"github.com/imgee/imgee/config"
	"github.com/imgee/imgee/httpd/mux"
	"github.com/imgee/imgee/httpd/services"
)

// init http handler
func Init() {
	mux.Handle("/api/resize", services.HandleResize)
}

// run http server by config
func Run() {
	var (
		endRunning = make(chan bool, 1)
		err        error
	)

	// http server
	if config.Conf.Web.EnableHttp {
		go func() {
			println("http server Running on ", config.Conf.Web.HttpPort)
			err = http.ListenAndServe(fmt.Sprintf(":%d", config.Conf.Web.HttpPort), mux.Muxd)
			if err != nil {
				println("ListenAndServe: ", err)
				endRunning <- true
			}
		}()
	}

	// https server
	if config.Conf.Web.EnableHttps {
		go func() {
			println("https server Running on ", config.Conf.Web.HttpsPort)
			err = http.ListenAndServeTLS(fmt.Sprintf(":%d", config.Conf.Web.HttpsPort), config.Conf.Web.CertFile, config.Conf.Web.KeyFile, mux.Muxd)
			if err != nil {
				println("ListenAndServeTLS: ", err)
				endRunning <- true
			}
		}()
	}
	<-endRunning
}
