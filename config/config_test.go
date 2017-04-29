// Package config provides ...
package config

import "testing"

// print Conf
func TestInit(t *testing.T) {
	Init("v0.0.1")
	t.Log(Conf)
}
