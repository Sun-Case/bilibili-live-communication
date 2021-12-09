package biliJsonConv

import (
	jsoniter "github.com/json-iterator/go"
)

type onlineRankV2 struct {
	Cmd  string `json:"cmd"`
	Data struct {
		List []struct {
			UId        jsoniter.Number `jsoniter:"uid"`
			Face       string          `jsoniter:"face"`
			Score      string          `jsoniter:"score"`
			UName      string          `jsoniter:"uname"`
			Rank       jsoniter.Number `jsoniter:"rank"`
			GuardLevel jsoniter.Number `jsoniter:"guard_level"`
		} `jsoniter:"list"`
		RankType string `jsoniter:"rank_type"`
	} `jsoniter:"data"`
}
type OnlineRankV2Struct struct {
	RawData onlineRankV2

	Rank []struct {
		UId        string
		Face       string
		Score      string
		UName      string
		Rank       string
		GuardLevel string
	}
}

func OnlineRankV2(b []byte) (OnlineRankV2Struct, error) {
	var olr onlineRankV2
	var olrs OnlineRankV2Struct
	var err error

	if err = jsoniter.Unmarshal(b, &olr); err != nil {
		goto Label
	}

	olrs.RawData = olr
	for i := range olr.Data.List {
		olrs.Rank = append(olrs.Rank, struct {
			UId        string
			Face       string
			Score      string
			UName      string
			Rank       string
			GuardLevel string
		}{UId: olr.Data.List[i].UId.String(),
			Face:       olr.Data.List[i].Face,
			Score:      olr.Data.List[i].Score,
			UName:      olr.Data.List[i].UName,
			Rank:       olr.Data.List[i].Rank.String(),
			GuardLevel: olr.Data.List[i].GuardLevel.String()})
	}

Label:
	return olrs, err
}
