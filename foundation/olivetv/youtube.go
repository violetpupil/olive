package olivetv

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/html"
)

func init() {
	registerSite("youtube", &youtube{})
}

type youtube struct {
	base
}

func (this *youtube) Name() string {
	return "油管"
}

func (this *youtube) Snap(tv *TV) error {
	tv.Info = &Info{
		Timestamp: time.Now().Unix(),
	}

	return this.set(tv)
}

func (this *youtube) set(tv *TV) error {
	liveURL := fmt.Sprintf("https://www.youtube.com/%s/live", tv.RoomID)

	// resp, err := req.C().R().Get(liveURL)
	// if err != nil {
	// 	return fmt.Errorf("get video url failed: %w", err)
	// }
	// videoHTML := resp.Bytes()

	resp, err := http.Get(liveURL)
	if err != nil {
		return err
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	videoHTML := content

	// PlayerResponse
	prb := this.playerResponseBytes(videoHTML)
	if len(prb) == 0 {
		// return fmt.Errorf("unable to retrieve player response object from watch page")
		return nil
	}
	var pr YoutubePlayerResponse

	err = jsoniter.Unmarshal(prb, &pr)
	if err != nil {
		return fmt.Errorf("unmarshal pr failed: %w", err)
	}

	tv.streamURL = fmt.Sprintf("https://www.youtube.com/watch?v=%s", pr.VideoDetails.VideoID)
	tv.roomOn = pr.VideoDetails.IsLive
	tv.roomName = pr.VideoDetails.Title
	tv.streamerName = pr.VideoDetails.Author

	return nil
}

// playerResponseBytes searches the given HTML for the player response object
func (this *youtube) playerResponseBytes(data []byte) []byte {
	playerRespDecl := []byte("var ytInitialPlayerResponse =")
	var objData []byte
	reader := bytes.NewReader(data)
	tokenizer := html.NewTokenizer(reader)
	isScript := false

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return objData
		case html.TextToken:
			if isScript {
				data := tokenizer.Text()
				declStart := bytes.Index(data, playerRespDecl)
				if declStart < 0 {
					continue
				}

				// Maybe add a LogTrace in the future for stuff like this
				//LogDebug("Found script element with player response in watch page.")
				objStart := bytes.Index(data[declStart:], []byte("{")) + declStart
				objEnd := bytes.LastIndex(data[objStart:], []byte("};")) + 1 + objStart

				if objEnd > objStart {
					objData = data[objStart:objEnd]
				}

				return objData
			}
		case html.StartTagToken:
			tn, _ := tokenizer.TagName()
			if string(tn) == "script" {
				isScript = true
			} else {
				isScript = false
			}
		}
	}
}

type YoutubePlayerResponse struct {
	ResponseContext struct {
		MainAppWebResponseContext struct {
			LoggedOut bool `json:"loggedOut"`
		} `json:"mainAppWebResponseContext"`
	} `json:"responseContext"`
	PlayabilityStatus struct {
		Status            string `json:"status"`
		Reason            string `json:"reason"`
		LiveStreamability struct {
			LiveStreamabilityRenderer struct {
				VideoID      string `json:"videoId"`
				OfflineSlate struct {
					LiveStreamOfflineSlateRenderer struct {
						ScheduledStartTime string `json:"scheduledStartTime"`
					} `json:"liveStreamOfflineSlateRenderer"`
				} `json:"offlineSlate"`
				PollDelayMs string `json:"pollDelayMs"`
			} `json:"liveStreamabilityRenderer"`
		} `json:"liveStreamability"`
	} `json:"playabilityStatus"`
	StreamingData struct {
		ExpiresInSeconds string `json:"expiresInSeconds"`
		AdaptiveFormats  []struct {
			Itag              int     `json:"itag"`
			URL               string  `json:"url"`
			MimeType          string  `json:"mimeType"`
			QualityLabel      string  `json:"qualityLabel,omitempty"`
			TargetDurationSec float64 `json:"targetDurationSec"`
		} `json:"adaptiveFormats"`
		DashManifestURL string `json:"dashManifestUrl"`
		HlsManifestURL  string `json:"hlsManifestUrl"`
	} `json:"streamingData"`
	VideoDetails struct {
		VideoID          string  `json:"videoId"`
		Title            string  `json:"title"`
		LengthSeconds    string  `json:"lengthSeconds"`
		IsLive           bool    `json:"isLive"`
		ChannelID        string  `json:"channelId"`
		IsOwnerViewing   bool    `json:"isOwnerViewing"`
		ShortDescription string  `json:"shortDescription"`
		AverageRating    float64 `json:"averageRating"`
		AllowRatings     bool    `json:"allowRatings"`
		ViewCount        string  `json:"viewCount"`
		Author           string  `json:"author"`
		IsLiveContent    bool    `json:"isLiveContent"`
	} `json:"videoDetails"`
	Microformat struct {
		PlayerMicroformatRenderer struct {
			Thumbnail struct {
				Thumbnails []struct {
					URL string `json:"url"`
				} `json:"thumbnails"`
			} `json:"thumbnail"`
			LiveBroadcastDetails struct {
				IsLiveNow      bool   `json:"isLiveNow"`
				StartTimestamp string `json:"startTimestamp"`
				EndTimestamp   string `json:"endTimestamp"`
			} `json:"liveBroadcastDetails"`
			PublishDate string `json:"publishDate"`
			UploadDate  string `json:"uploadDate"`
		} `json:"playerMicroformatRenderer"`
	} `json:"microformat"`
}
