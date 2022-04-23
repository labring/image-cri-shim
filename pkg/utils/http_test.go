/*
Copyright 2022 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestHTTP(t *testing.T) {
	//192.168.5.58:5000
	prefix := "http://192.168.5.58:5000"
	//byte[] authEncBytes = Base64.encodeBase64(authString.getBytes("utf-8"));
	type RegistryData struct {
		Name string
		Tags []string
	}
	var registry RegistryData
	data, _ := HTTP(fmt.Sprintf("%s/v2/sealyun/lvscare/tags/list", prefix), map[string]string{"Authorization": "Basic YWRtaW46cGFzc3cwcmQ="})
	if data != "" {
		json.Unmarshal([]byte(data), &registry)
		t.Log(data)
	}
}
