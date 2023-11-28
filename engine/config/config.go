// Package config stores config type.
package config

import (
	"os"
	"os/exec"
	"time"

	"github.com/imdario/mergo"
)

const CoreConfigKey = "core_config"

var DefaultConfig = Config{
	PortalUsername: "olive",
	PortalPassword: "olive",

	// core
	LogDir:                   "",
	SaveDir:                  "",
	OutTmpl:                  `[{{ .StreamerName }}][{{ .RoomName }}][{{ now | date "2006-01-02 15-04-05"}}].flv`,
	LogLevel:                 5,
	SnapRestSeconds:          15,
	SplitRestSeconds:         10,
	CommanderPoolSize:        1,
	ParserMonitorRestSeconds: 300,

	// tv
	DouyinCookie:   "",
	KuaishouCookie: "",
	TikTokProxy:    "",

	PostCmdsRetryCount: 2,
}

type Config struct {
	// portal
	PortalUsername string
	PortalPassword string

	// core
	LogDir                   string
	SaveDir                  string
	OutTmpl                  string
	LogLevel                 uint32
	SnapRestSeconds          uint
	SplitRestSeconds         uint
	CommanderPoolSize        uint
	ParserMonitorRestSeconds uint
	PostCmdsRetryCount       int

	// tv
	DouyinCookie   string
	KuaishouCookie string
	TikTokProxy    string

	// biliup
	BiliupEnable      bool
	CookieFilepath    string
	Threads           int64
	MaxBytesPerSecond float64
}

func (cfg *Config) CheckAndFix() {
	wd, _ := os.Getwd()
	DefaultConfig.LogDir = wd
	DefaultConfig.SaveDir = wd

	mergo.Merge(cfg, DefaultConfig)
}

type ID string

type Bout interface {
	IsConfigValid() bool
	// show settings
	GetID() ID
	GetPlatform() string
	GetRoomID() string
	GetStreamerName() string
	GetOutFilename() string
	GetOutTmpl() string
	GetSaveDir() string
	GetParser() string
	GetPostCmds() []*exec.Cmd
	SatisfySplitRule(time.Time, string) bool

	// show events
	AddMonitor() error
	RemoveMonitor() error
	AddRecorder() error
	RemoveRecorder() error
	RestartRecorder()

	// tv
	Snap() error
	StreamURL() (string, bool)
	RoomName() (string, bool)
	StreamerName() (string, bool)
	SiteName() string
}
