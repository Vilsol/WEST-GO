package WEST

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type WESTRequest struct {
	Id     string
	Method string
	Path   string
	Body   string
}

func processRequest(message []byte) (*WESTRequest, error) {
	length := len(message)

	i := 0

	i, id := readUntilSpace(i, message)

	if len(id) < 1 || len(id) > 10 {
		return nil, errors.New("Id must be between 1 and 10 bytes")
	}

	i, method := readUntilSpace(i, message)

	if !stringInSlice(method, methods) {
		return &WESTRequest{
			Id: id,
		}, errors.New("Method must be one of: " + strings.Join(methods, ", "))
	}

	pathLength := length - i

	if method == "POST" || method == "PUT" || method == "PATCH" {
		temp, pathLengthString := readUntilSpace(i, message)

		i = temp

		pathLength, _ = strconv.Atoi(pathLengthString)
	}

	path := string(message[i : i+pathLength])

	if len(path) < 1 {
		return &WESTRequest{
			Id: id,
		}, errors.New("Invalid path")
	}

	i += pathLength

	body := string(message[i:])

	return &WESTRequest{
		Id:     id,
		Method: method,
		Path:   path,
		Body:   body,
	}, nil
}

func readUntilSpace(start int, data []byte) (int, string) {
	read := ""
	i := start

	for i = start; i < len(data); i++ {
		if data[i] == ' ' {
			i++
			break
		}

		read += string(data[i])
	}

	return i, read
}

func (west WESTRequest) toHTTPRequest(originalRequest *http.Request) *http.Request {
	request := http.Request{
		Method:     west.Method,
		RequestURI: west.Path,
		RemoteAddr: originalRequest.RemoteAddr,
		Host:       originalRequest.Host,
		URL: &url.URL{
			Host: originalRequest.Host,
			Path: west.Path,
		},
	}

	if len(west.Body) > 0 {
		request.Body = ioutil.NopCloser(bytes.NewReader([]byte(west.Body)))
	}

	return &request
}
