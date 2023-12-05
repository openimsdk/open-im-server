// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package body

import (
	"errors"

	"github.com/OpenIMSDK/protocol/constant"
)

const (
	ANDROID      = "android"
	IOS          = "ios"
	QUICKAPP     = "quickapp"
	WINDOWSPHONE = "winphone"
	ALL          = "all"
)

type Platform struct {
	Os     any
	osArry []string
}

func (p *Platform) Set(os string) error {
	if p.Os == nil {
		p.osArry = make([]string, 0, 4)
	} else {
		switch p.Os.(type) {
		case string:
			return errors.New("platform is all")
		default:
		}
	}

	for _, value := range p.osArry {
		if os == value {
			return nil
		}
	}

	switch os {
	case IOS:
		fallthrough
	case ANDROID:
		fallthrough
	case QUICKAPP:
		fallthrough
	case WINDOWSPHONE:
		p.osArry = append(p.osArry, os)
		p.Os = p.osArry
	default:
		return errors.New("unknow platform")
	}

	return nil
}

func (p *Platform) SetPlatform(platform string) error {
	switch platform {
	case constant.AndroidPlatformStr:
		return p.SetAndroid()
	case constant.IOSPlatformStr:
		return p.SetIOS()
	default:
		return errors.New("platform err")
	}
}

func (p *Platform) SetIOS() error {
	return p.Set(IOS)
}

func (p *Platform) SetAndroid() error {
	return p.Set(ANDROID)
}

func (p *Platform) SetQuickApp() error {
	return p.Set(QUICKAPP)
}

func (p *Platform) SetWindowsPhone() error {
	return p.Set(WINDOWSPHONE)
}

func (p *Platform) SetAll() {
	p.Os = ALL
}
