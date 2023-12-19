package configs

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Setting struct {
	vp *viper.Viper
}

func NewSetting() (*Setting, error) {
	vp := viper.New()
	vp.SetConfigType("yaml")
	vp.SetConfigName("config")

	if true {
		isHome := viper.Get("is_home")
		if isHome.(bool) {
			vp.AddConfigPath("user-web/configs/home")
		} else {
			vp.AddConfigPath("user-web/configs/debug")
		}
	} else {
		vp.AddConfigPath("user-web/configs/production")
	}

	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}

	s := &Setting{vp}
	s.WatchSettingChange()
	return s, nil
}

var sections = make(map[string]any)

func (s *Setting) ReadSection(k string, v any) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return nil
	}
	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}

func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		if err := s.ReadSection(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(in fsnotify.Event) {
			_ = s.ReloadAllSection()
		})
	}()
}
