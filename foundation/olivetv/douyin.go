package olivetv

import (
	"errors"
	"strings"
	"time"

	"github.com/go-olive/olive/foundation/olivetv/model"
	"github.com/imroc/req/v3"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

var (
	ErrCookieNotSet = errors.New("cookie not configured")
)

func init() {
	registerSite("douyin", &douyin{})
}

type douyin struct {
	base
}

func (this *douyin) Name() string {
	return "抖音"
}

func (this *douyin) Snap(tv *TV) error {
	tv.Info = &Info{
		Timestamp: time.Now().Unix(),
	}
	return this.set(tv)
}

const CHROME = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36"

func (this *douyin) set(tv *TV) error {
	const douyincookie = `ttwid=1%7CcfLDfqNkk8o-IKEppAbVXFkCcglSlTQbXQ-sOIZqeT0%7C1693836521%7C5a64077ccdaa38c03827e07363efbb814bd9000c50c9e3247146015388928089; home_can_add_dy_2_desktop=%220%22; strategyABtestKey=%221693836522.123%22; stream_recommend_feed_params=%22%7B%5C%22cookie_enabled%5C%22%3Atrue%2C%5C%22screen_width%5C%22%3A1920%2C%5C%22screen_height%5C%22%3A1080%2C%5C%22browser_online%5C%22%3Atrue%2C%5C%22cpu_core_num%5C%22%3A8%2C%5C%22device_memory%5C%22%3A8%2C%5C%22downlink%5C%22%3A10%2C%5C%22effective_type%5C%22%3A%5C%224g%5C%22%2C%5C%22round_trip_time%5C%22%3A50%7D%22; FORCE_LOGIN=%7B%22videoConsumedRemainSeconds%22%3A180%7D; volume_info=%7B%22isUserMute%22%3Afalse%2C%22isMute%22%3Afalse%2C%22volume%22%3A0.5%7D; passport_csrf_token=c5b117b9f262acf5d07f6a29b0425ab4; passport_csrf_token_default=c5b117b9f262acf5d07f6a29b0425ab4; odin_tt=d490a5b68981828937351d90d52d717fdcacd3f038f0f0505efb8a203d40e9744b5fbefe9c259879d1bffcc91a39d608a958d6584b53bb2fa9b6cc07b10777baa725728e6fa2f25940f68182b681f731; VIDEO_FILTER_MEMO_SELECT=%7B%22expireTime%22%3A1694441331554%2C%22type%22%3A1%7D; IsDouyinActive=false; msToken=ctKJrLjrxRieZ2vWcv0ZCTcFAtLdtAMSSSBw85YkTNegSzakztz0H8UIiJSAZ0CKv1gq5EqxVUe8lt1XDBKjpQ-JGwfhNpN3iVyb1ntT45xSYT5c6g==`
	if tv.cookie == "" || strings.Contains(tv.cookie, "ac_nonce") {
		tv.cookie = douyincookie
	}
	api := `https://live.douyin.com/webcast/room/web/enter/?aid=6383&live_id=1&device_platform=web&language=zh-CN&enter_from=web_live&cookie_enabled=true&screen_width=1536&screen_height=864&browser_language=zh-CN&browser_platform=Win32&browser_name=Chrome&browser_version=94.0.4606.81&room_id_str=&enter_source=&web_rid=` +
		tv.RoomID
	resp, err := req.C().R().
		SetHeaders(map[string]string{
			"User-Agent":      CHROME,
			"referer":         "https://live.douyin.com/",
			"cookie":          tv.cookie,
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		}).
		Get(api)
	if err != nil {
		return err
	}
	// log.Println(api)
	text := resp.String()
	text = gjson.Get(text, "data.data.0").String()
	// 抖音 status == 2 代表是开播的状态
	if gjson.Get(text, "status").String() != "2" {
		return nil
	}

	streamDataStr := gjson.Get(text, "stream_url.live_core_sdk_data.pull_data.stream_data").String()
	var streamData model.DouyinStreamData
	err = jsoniter.UnmarshalFromString(streamDataStr, &streamData)
	if err != nil {
		return err
	}
	flv := streamData.Data.Origin.Main.Flv
	tv.streamURL = flv
	tv.roomOn = true
	tv.roomName = gjson.Get(text, "title").String()

	return nil
}
