package controllers

import (
	"net/http"
	"io/ioutil"
	"strings"
	"fmt"
	"github.com/dugancathal/stuffs/config"
	"encoding/json"
)

type EchoController struct {
	endpointStorage config.RouteMapping
	requestData map[string][][]byte
}

func NewEchoController() *EchoController {
	return NewEchoControllerFromConfig(make(config.RouteMapping))
}

func NewEchoControllerFromConfig(initialRouteMapping config.RouteMapping) *EchoController {
	return &EchoController{endpointStorage: initialRouteMapping, requestData: make(map[string][][]byte)}
}

func (self *EchoController) HandleReq(writer http.ResponseWriter, request *http.Request) {
	endpoint := request.URL.Path
	method := request.Method

	requestBody, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	if _, ok := self.requestData[fmt.Sprintf("%s %s", method, endpoint)]; !ok {
		self.requestData[fmt.Sprintf("%s %s", method, endpoint)] = [][]byte{}
	}
	self.requestData[fmt.Sprintf("%s %s", method, endpoint)] = append(self.requestData[fmt.Sprintf("%s %s", method, endpoint)], requestBody)
	fmt.Println(self.requestData)
	body, _ := self.endpointStorage[fmt.Sprintf("%s %s", method, endpoint)]
	writer.Write(body)
}

func (self *EchoController) HandleGetReq(writer http.ResponseWriter, request *http.Request) {
	endpoint := parseEndpoint(request)
	method := request.Method

	bodies, _ := self.requestData[fmt.Sprintf("%s %s", method, endpoint)]
	body, _ := json.Marshal(bodies)
	writer.Write(body)
}

func (self *EchoController) HandleSetReq(writer http.ResponseWriter, request *http.Request) {
	endpoint := parseEndpoint(request)
	method := request.Method
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return
	}

	self.endpointStorage[fmt.Sprintf("%s %s", method, endpoint)] = body
}

func (self *EchoController) GetConfiguration(writer http.ResponseWriter, request *http.Request) {
	prettyMap := make(map[string]string, len(self.endpointStorage))
	for k, v := range self.endpointStorage {
		prettyMap[k] = string(v)
	}
	body, _ := json.Marshal(prettyMap)
	writer.Write(body)
}

func parseEndpoint(request *http.Request) string {
	parts := strings.Split(request.URL.Path, "/")
	endpoint := "/" + strings.Join(parts[2:], "/")
	return endpoint
}
