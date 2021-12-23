package main

import (
	"bilibili-live-communication/bilibili/biliAPI"
	"bilibili-live-communication/bilibili/biliDataConv"
	"bilibili-live-communication/bilibili/biliJsonConv"
	"bilibili-live-communication/bilibili/biliWebsocket"
	"fmt"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"log"
	"os"
	"strconv"
	"time"
)

var KV map[string]string
var show = map[biliJsonConv.CmdType]bool{
	biliJsonConv.DanMuMsgType:                  false,
	biliJsonConv.InteractWordType:              false,
	biliJsonConv.OnlineRankCountType:           false,
	biliJsonConv.StopLiveRoomListType:          false,
	biliJsonConv.RoomRealTimeMessageUpdateType: false,
	biliJsonConv.OnlineRankV2Type:              false,
	biliJsonConv.HotRankChangedV2Type:          false,
	biliJsonConv.PreparingType:                 false,
	biliJsonConv.SendGiftType:                  false,
	biliJsonConv.EntryEffectType:               false,
	biliJsonConv.LiveInteractiveGameType:       false,
	biliJsonConv.ComboSendType:                 false,
	biliJsonConv.NoticeMsgType:                 false,
	biliJsonConv.OnlineRankTop3Type:            false,
	biliJsonConv.LiveType:                      false,
	biliJsonConv.RoomChangeType:                false,
	biliJsonConv.WidgetBannerType:              false,
	biliJsonConv.GuardBuyType:                  true,
	biliJsonConv.ErrorType:                     true,
	biliJsonConv.UnknownType:                   true,
}

func main() {
	fb, _ := os.ReadFile("config.json")

	_ = jsoniter.Unmarshal(fb, &KV)

	rid := KV["room_id"]

	// 获取真正直播间ID
	rid, err := biliAPI.GetRoomInfo(rid)
	if err != nil {
		log.Panicln(err)
	}

	// 获取弹幕服务器地址及token
	hosts, token, err := biliAPI.GetDanmuInfo(rid)
	if err != nil {
		log.Panicln(err)
	}

	wsHost := fmt.Sprintf("wss://%s/sub", hosts[0])
	log.Println(wsHost)
	c, _, err := websocket.DefaultDialer.Dial(wsHost, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Println(err)
		}
	}(c)

	sendCannel := make(chan biliDataConv.EncodeArgs)
	receiveChannel := make(chan []byte)
	blw := biliWebsocket.New(c, &receiveChannel, &sendCannel, func(liveWebsocket *biliWebsocket.BiliWebsocket) biliDataConv.EncodeArgs {
		// 注意，如果获取弹幕地址及Token时，没有携带账号的Cookie，则下面的uid必须为0，否则鉴权失败，直接断开连接
		data := fmt.Sprintf("{\"uid\":0,\"roomid\":%s,\"protover\":3,\"platform\":\"web\",\"type\":2,\"key\":\"%s\"}", rid, token)
		return biliDataConv.EncodeArgs{
			Version: 0x0001, Operation: 0x00000007, Sequence: 0x00000001, Body: []byte(data),
		}
	})

	for {
		select {
		case <-blw.IsClosedCh:
			return
		case b, ok := <-receiveChannel:
			if !ok {
				return
			}
			switch biliJsonConv.GetCmdType(b) {
			case biliJsonConv.DanMuMsgType:
				if b, err := biliJsonConv.DanMuMsg(b); err == nil {
					if show[biliJsonConv.DanMuMsgType] {
						log.Printf("%s: %s", b.Name, b.Msg)
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.InteractWordType:
				if b, err := biliJsonConv.InteractWord(b); err == nil {
					if show[biliJsonConv.InteractWordType] {
						log.Println(b.UName, "进入直播间")
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.OnlineRankCountType:
				if b, err := biliJsonConv.OnlineRankCount(b); err == nil {
					if show[biliJsonConv.OnlineRankCountType] {
						log.Println("在线排名统计:", b.Count)
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.StopLiveRoomListType:
				if b, err := biliJsonConv.StopLiveRoomList(b); err == nil {
					if show[biliJsonConv.StopLiveRoomListType] {
						log.Println("下播直播间ID:", b.RoomIdList)
					}
				}
			case biliJsonConv.RoomRealTimeMessageUpdateType:
				if b, err := biliJsonConv.RoomRealTimeMessageUpdate(b); err == nil {
					if show[biliJsonConv.RoomRealTimeMessageUpdateType] {
						log.Printf("直播间: %s, 粉丝数: %s, 粉丝团数: %s\n", b.RoomId, b.Fans, b.FansClub)
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.OnlineRankV2Type:
				if b, err := biliJsonConv.OnlineRankV2(b); err == nil {
					if show[biliJsonConv.OnlineRankV2Type] {
						log.Println("高能榜: ")
						for i := range b.Rank {
							fmt.Printf("%s. %s    ", b.Rank[i].Rank, b.Rank[i].UName)
						}
						fmt.Println("")
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.HotRankChangedV2Type:
				if b, err := biliJsonConv.HotRankChangedV2(b); err == nil {
					if show[biliJsonConv.HotRankChangedV2Type] {
						log.Printf("%s排行榜 刷新: %s 之 %s\n", b.AreaName, b.RankDesc, b.Rank)
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.PreparingType:
				if b, err := biliJsonConv.Preparing(b); err == nil {
					if show[biliJsonConv.PreparingType] {
						log.Printf("直播间: %s 下播\n", b.RoomId)
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.SendGiftType:
				if b, err := biliJsonConv.SendGift(b); err == nil {
					if show[biliJsonConv.SendGiftType] {
						log.Println(b.Uname, b.Action, b.GiftName)
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.EntryEffectType:
				if b, err := biliJsonConv.EntryEffect(b); err == nil {
					if show[biliJsonConv.EntryEffectType] {
						log.Println(b.CopyWritingV2)
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.LiveInteractiveGameType:
				if b, err := biliJsonConv.LiveInteractiveGame(b); err == nil {
					if show[biliJsonConv.LiveInteractiveGameType] {
						log.Println(b.Uname, "喂投", b.GiftName, "共", b.GiftNum, "个")
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.ComboSendType:
				if b, err := biliJsonConv.ComboSend(b); err == nil {
					if show[biliJsonConv.ComboSendType] {
						log.Println(b.Uname, b.Action, b.GiftName, "共", b.ComboNum, "个")
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.NoticeMsgType:
				if b, err := biliJsonConv.NoticeMsg(b); err == nil {
					if show[biliJsonConv.NoticeMsgType] {
						log.Println(b.MsgCommon)
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.OnlineRankTop3Type:
				if b, err := biliJsonConv.OnlineRankTop3(b); err == nil {
					if show[biliJsonConv.OnlineRankTop3Type] {
						for i := range b.List {
							log.Println(b.List[i].Msg, "第", b.List[i].Rank)
						}
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.LiveType:
				if b, err := biliJsonConv.Live(b); err == nil {
					if show[biliJsonConv.LiveType] {
						log.Println(b.RoomId, "开播", "设备:", b.LivePlatform)
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.RoomChangeType:
				if b, err := biliJsonConv.RoomChange(b); err == nil {
					if show[biliJsonConv.RoomChangeType] {
						log.Println("直播间:", b.Title, ", 总分区:", b.ParentAreaName, "子分区:", b.AreaName)
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.WidgetBannerType:
				if b, err := biliJsonConv.WidgetBanner(b); err == nil {
					if show[biliJsonConv.WidgetBannerType] {
						log.Println(b)
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.GuardBuyType:
				if b, err := biliJsonConv.GuardBuy(b); err == nil {
					if show[biliJsonConv.GuardBuyType] {
						log.Println(b.Username, b.GiftName, b.Num)
					}
				} else {
					log.Println(err)
				}
			case biliJsonConv.ErrorType:
				log.Println("Error Type:", string(b))
			case biliJsonConv.UnknownType:
				fallthrough
			default:
				log.Println(string(b))
				save(b)
			}
		}
	}
}

func save(b []byte) {
	var cmd struct {
		Cmd string `json:"cmd"`
	}
	if err := jsoniter.Unmarshal(b, &cmd); err != nil {
		return
	} else {
		if err := os.MkdirAll(KV["save_path"]+cmd.Cmd, 0755); err != nil {
			return
		} else {
			_ = os.WriteFile(KV["save_path"]+cmd.Cmd+"/"+strconv.FormatUint(uint64(time.Now().UnixNano()), 10)+".json", b, 0644)
		}
	}
}
