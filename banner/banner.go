// Package banner provides ...
package banner

import (
	"fmt"
)

const b = `
=================================================
                                     _
           ___  ___  _ __   ___ __ _(_)
          / _ \/ _ \| '_ \ / _ ' _' | |
          \__  \__  | |_) | | | | | | |
          |___/|___/| .__/|_| |_| |_|_|
                     \___|             

               Image Editor Edit
       Opensource image processing tools
                http://imgee.me
================== Imgee %s ==================
	`

// Println print out banner information
func Println(version string) {
	fmt.Printf("\033[2J")           // clear screen
	fmt.Printf("\033[%d;%dH", 0, 0) // move cursor to x-0, y=0
	fmt.Printf(b, version)          // print banner including version
}
