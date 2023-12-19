package configs

import "time"

type ServerSettingS struct {
	RunMode      string
	HttpHost     string
	HttpPort     string
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	RemoteConfig RemoteServerSettingS
}

type RemoteServerSettingS struct {
	Host string
	Port string
}
