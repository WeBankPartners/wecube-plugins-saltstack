package test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

var resourceIds = make(map[string]string)

const (
	PLUGIN_HOST_URL = "http://localhost:8082"
)

type Outputs struct {
	Outputs []Output `json:"outputs,omitempty"`
}

type Output struct {
	Guid string `json:"guid,omitempty"`
}

type PluginResponse struct {
	ResultCode string  `json:"result_code"`
	ResultMsg  string  `json:"result_message"`
	Results    Outputs `json:"results"`
}

func CallPlugin(name, action, input string) map[string]string {
	output, err := http.Post(PLUGIN_HOST_URL+"/v1/deploy/"+name+"/"+action, "application/json", strings.NewReader(input))

	if err != nil {
		logrus.Errorf("call plugin server meet error = %v", err)
	}

	pluginResponse := PluginResponse{}
	err = UnmarshalJson(output.Body, &pluginResponse)
	if err != nil {
		logrus.Errorf("unmarshal plugin response meet error = %v", err)
	}

	if pluginResponse.ResultCode == "1" {
		logrus.Errorf("call plugin meet error = %v", pluginResponse.ResultMsg)
	}

	outputMap := make(map[string]string)
	for i := 0; i < len(pluginResponse.Results.Outputs); i++ {
		outputMap[pluginResponse.Results.Outputs[i].Guid] = pluginResponse.Results.Outputs[i].Guid
	}
	logrus.Infof("resource (ids=%v) have been handled, plugin = %v, action = %v", outputMap, name, action)

	return outputMap
}

func UnmarshalJson(source interface{}, target interface{}) error {
	reader, ok := source.(io.Reader)
	if !ok {
		return fmt.Errorf("the source to be unmarshaled is not a io.reader type")
	}

	bodyBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("parse http request (%v) meet error (%v)", reader, err)
	}

	logrus.Infof("http response = %v", string(bodyBytes))

	if err = json.Unmarshal(bodyBytes, target); err != nil {
		return fmt.Errorf("unmarshal http request (%v) meet error (%v)", reader, err)
	}

	return nil
}
