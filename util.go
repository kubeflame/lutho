package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

const END_OF_TRANSMISSION = "\u0004"

var WSCloseCode = struct {
	info, warning, error uint
}{
	info:    3000,
	warning: 3001,
	error:   3002,
}

var WSOpType = struct {
	bind,
	selfSubjectAccessReview,
	list,
	listAll,
	helmList,
	get,
	helmShowValues,
	helmGet,
	helmInstall,
	helmUpgrade,
	helmPull,
	helmGetTags,
	update,
	check,
	delete,
	helmUninstall,
	close,
	stdin,
	stdout,
	resize,
	toast string
}{
	bind:                    "bind",
	selfSubjectAccessReview: "selfSubjectAccessReview",
	list:                    "list",
	listAll:                 "listAll",
	helmList:                "helmList",
	get:                     "get",
	helmShowValues:          "helmShowValues",
	helmGet:                 "helmGet",
	helmInstall:             "helmInstall",
	helmUpgrade:             "helmUpgrade",
	helmPull:                "helmPull",
	helmGetTags:             "helmGetTags",
	check:                   "check",
	update:                  "update",
	delete:                  "delete",
	helmUninstall:           "helmUninstall",
	close:                   "close",
	stdin:                   "stdin",
	stdout:                  "stdout",
	resize:                  "resize",
	toast:                   "toast",
}

var KubernetesConfigType = struct {
	inClusterConfig,
	kubeconfigPath,
	kubeconfigRaw,
	accessToken string
}{
	inClusterConfig: "inClusterConfig",
	kubeconfigPath:  "kubeconfigPath",
	kubeconfigRaw:   "kubeconfigRaw",
	accessToken:     "accessToken",
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

func reverseStringSlice(input []string) {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
}

var errCouldNotFindResource = errors.New("the server could not find the requested resource")
