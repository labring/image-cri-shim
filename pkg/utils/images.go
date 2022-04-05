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

import "github.com/pkg/errors"

func LoadImages(imageDir string) ([]string, error) {
	var imageList []string
	if imageDir != "" && IsExist(imageDir) {
		paths, err := GetFiles(imageDir)

		if err != nil {
			return nil, errors.Wrap(err, "load image list files error")
		}
		for _, p := range paths {
			images, err := ReadLines(p)
			if err != nil {
				return nil, errors.Wrap(err, "load image list error")
			}
			imageList = append(imageList, images...)
		}
	}
	imageList = RemoveDuplicate(imageList)
	return imageList, nil
}
