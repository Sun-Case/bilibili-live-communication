package biliJsonConv

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
)

type entryEffect struct {
	Cmd  string `json:"cmd"`
	Data struct {
		Id               jsoniter.Number `json:"id"`
		Uid              jsoniter.Number `json:"uid"`
		TargetId         jsoniter.Number `json:"target_id"`
		MockEffect       jsoniter.Number `json:"mock_effect"`
		Face             string          `json:"face"`
		PrivilegeType    jsoniter.Number `json:"privilege_type"`
		CopyWriting      string          `json:"copy_writing"`
		CopyColor        string          `json:"copy_color"`
		HighlightColor   string          `json:"highlight_color"`
		Priority         jsoniter.Number `json:"priority"`
		BasemapUrl       string          `json:"basemap_url"`
		ShowAvatar       jsoniter.Number `json:"show_avatar"`
		EffectiveTime    jsoniter.Number `json:"effective_time"`
		WebBasemapUrl    string          `json:"web_basemap_url"`
		WebEffectiveTime jsoniter.Number `json:"web_effective_time"`
		WebEffectClose   jsoniter.Number `json:"web_effect_close"`
		WebCloseTime     jsoniter.Number `json:"web_close_time"`
		Business         jsoniter.Number `json:"business"`
		CopyWritingV2    string          `json:"copy_writing_v2"`
		IconList         []interface{}   `json:"icon_list"`
		MaxDelayTime     jsoniter.Number `json:"max_delay_time"`
		TriggerTime      jsoniter.Number `json:"trigger_time"`
		Identities       jsoniter.Number `json:"identities"`
	} `json:"data"`
}
type EntryEffectStruct struct {
	RawData entryEffect

	Id, Uid, TargetId                     string
	CopyWriting, CopyColor, CopyWritingV2 string
	TriggerTime                           string
}

func EntryEffect(b []byte) (EntryEffectStruct, error) {
	var ee entryEffect
	var ees EntryEffectStruct
	var err error

	decoder := jsoniter.NewDecoder(bytes.NewReader(b))
	if err = decoder.Decode(&ee); err != nil {
		goto Label
	}

	ees.Id = ee.Data.Id.String()
	ees.Uid = ee.Data.Uid.String()
	ees.TargetId = ee.Data.TargetId.String()
	ees.CopyWriting = ee.Data.CopyWriting
	ees.CopyColor = ee.Data.CopyColor
	ees.CopyWritingV2 = ee.Data.CopyWritingV2
	ees.TriggerTime = ee.Data.TriggerTime.String()

Label:
	return ees, err
}
