package main

import (
	"crypto/rand"
	"encoding/hex"
)

const END_OF_TRANSMISSION = "\u0004"

var WSCloseCode = struct {
	info, warning, error uint32
}{
	info:    3000,
	warning: 3001,
	error:   3002,
}

var WSOpType = struct {
	bind,
	list,
	helmList,
	get,
	helmGet,
	update,
	check,
	delete,
	close,
	stdin,
	stdout,
	resize,
	toast string
}{
	bind:     "bind",
	list:     "list",
	helmList: "helmList",
	get:      "get",
	helmGet:  "helmGet",
	check:    "check",
	update:   "update",
	delete:   "delete",
	close:    "close",
	stdin:    "stdin",
	stdout:   "stdout",
	resize:   "resize",
	toast:    "toast",
}

func genSessionId() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	id := make([]byte, hex.EncodedLen(len(bytes)))
	hex.Encode(id, bytes)
	return string(id), nil
}

func isValidShell(validShells []string, shell string) bool {
	for _, validShell := range validShells {
		if validShell == shell {
			return true
		}
	}
	return false
}
