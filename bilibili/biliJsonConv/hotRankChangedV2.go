package biliJsonConv

import jsoniter "github.com/json-iterator/go"

type hotRankChangedV2 struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Rank        jsoniter.Number `json:"rank"`
		Trend       jsoniter.Number `json:"trend"`
		Countdown   jsoniter.Number `json:"countdown"`
		Timestamp   jsoniter.Number `json:"timestamp"`
		WebUrl      string          `json:"web_url"`
		LiveUrl     string          `json:"live_url"`
		BlinkUrl    string          `json:"blink_url"`
		LiveLinkUrl string          `json:"live_link_url"`
		PcLinkUrl   string          `json:"pc_link_url"`
		Icon        string          `json:"icon"`
		AreaName    string          `json:"area_name"`
		RankDesc    string          `json:"rank_desc"`
	} `json:"data"`
}
type HotRankChangedV2Struct struct {
	RawData hotRankChangedV2

	Timestamp string
	Rank      string
	AreaName  string
	RankDesc  string
}

func HotRankChangedV2(b []byte) (HotRankChangedV2Struct, error) {
	var hrc hotRankChangedV2
	var hrcs HotRankChangedV2Struct
	var err error

	if err = jsoniter.Unmarshal(b, &hrc); err != nil {
		goto Label
	}

	hrcs.RawData = hrc
	hrcs.Timestamp = hrc.Data.Timestamp.String()
	hrcs.Rank = hrc.Data.Rank.String()
	hrcs.AreaName = hrc.Data.AreaName
	hrcs.RankDesc = hrc.Data.RankDesc

Label:
	return hrcs, err
}
