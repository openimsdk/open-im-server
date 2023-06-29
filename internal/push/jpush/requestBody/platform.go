package requestBody

import (
	"Open_IM/pkg/common/constant"
	"errors"
)

const (
	ANDROID      = "android"
	IOS          = "ios"
	QUICKAPP     = "quickapp"
	WINDOWSPHONE = "winphone"
	ALL          = "all"
)

type Platform struct {
	Os     interface{}
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
