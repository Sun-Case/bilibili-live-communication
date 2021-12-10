package biliJsonConv

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
)

type onlineRankTop3 struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Dmscore jsoniter.Number `json:"dmscore"`
		List    []struct {
			Msg  string          `json:"msg"`
			Rank jsoniter.Number `json:"rank"`
		} `json:"list"`
	} `json:"data"`
}

type OnlineRankTop3Struct struct {
	RawData onlineRankTop3

	DmScore string
	List    []struct {
		Msg  string
		Rank string
	}
}

func OnlineRankTop3(b []byte) (OnlineRankTop3Struct, error) {
	var ort onlineRankTop3
	var orts OnlineRankTop3Struct
	var err error

	decoder := jsoniter.NewDecoder(bytes.NewReader(b))

	if err = decoder.Decode(&ort); err != nil {
		goto Label
	}

	orts.RawData = ort
	orts.DmScore = ort.Data.Dmscore.String()

	for i := range ort.Data.List {
		orts.List = append(orts.List, struct {
			Msg  string
			Rank string
		}{Msg: ort.Data.List[i].Msg, Rank: ort.Data.List[i].Rank.String()})
	}

Label:
	return orts, err
}
