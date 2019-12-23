package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sethvargo/go-password/password"
)

type pwgenHandler struct {
	length      int
	numDigits   int
	numSymbols  int
	noUpper     bool
	allowRepeat bool
}

func (h *pwgenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := password.Generate(h.length, h.numDigits, h.numSymbols, h.noUpper, h.allowRepeat)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%+v\n", res)
	return
}

func PwgenHandler(length int, numDigits int, numSymbols int, noUpper bool, allowRepeat bool) http.Handler {
	return &pwgenHandler{
		length: length, numDigits: numDigits, numSymbols: numSymbols, noUpper: noUpper, allowRepeat: allowRepeat,
	}
}
