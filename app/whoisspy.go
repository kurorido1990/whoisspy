package app

import (
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

var gen Generator
var node *snowflake.Node
var roomList sync.Map

func test(c *gin.Context) {
	data := new(IndexData)
	data.Title = "首頁"
	data.Content = "我的第一個首頁"
	c.HTML(http.StatusOK, "index.html", data)
}

func Run() {
	server := gin.Default()
	server.Use(cors.Default())
	server.LoadHTMLGlob("app/template/*")
	server.GET("/", test)
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(Status_OK, gin.H{
			"message": "pong",
		})
	})

	server.GET("/createRoom/:maxLimit", func(ctx *gin.Context) {
		maxLimit, _ := strconv.ParseInt(ctx.Params.ByName("maxLimit"), 10, 64)

		roomID := createRoom(int(maxLimit))

		ctx.JSON(Status_OK, gin.H{
			"roomID": roomID,
		})
	})

	server.GET("/addPlayer/:roomID/:name", addPlayer)
	server.GET("/getCard/:roomID/:playerID", getCard)

	server.GET("/monitor", func(ctx *gin.Context) {
		byteRoomList, _ := json.Marshal(roomList)
		ctx.JSON(http.StatusOK, struct {
			RoomList string
		}{
			RoomList: string(byteRoomList),
		})
	})
	server.GET("/room/:roomID/:playerID", gamePage)
	server.GET("/resetRoom/:roomID", func(ctx *gin.Context) {
		roomID := ctx.Params.ByName("roomID")

		if room := getRoom(roomID); room != nil {
			room.resetGame()

			ctx.JSON(Status_OK, nil)
			return
		}

		ctx.JSON(400, "不明原因壞了")
	})

	server.GET("/kick/:roomID/:playerID", kickPlayer)

	initRoomList()
	initSnowflake()
	initGen()
	//server.Run(":9999")
	server.Run()
}

func initSnowflake() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

func initRoomList() {
	roomList = sync.Map{}
}

func initGen() {
	gen = CreateGen()
}

func kickPlayer(ctx *gin.Context) {
	roomID := ctx.Params.ByName("roomID")
	playerID := ctx.Params.ByName("playerID")

	if room := getRoom(roomID); room != nil {
		room.kickPlayer(playerID)
	}

	ctx.JSON(Status_OK, nil)
}

func gamePage(ctx *gin.Context) {
	roomID := ctx.Params.ByName("roomID")
	playerID := ctx.Params.ByName("playerID")

	if room := getRoom(roomID); room != nil {
		if room.Status == RoomStatusEnd {
			ctx.HTML(http.StatusOK, "end.html", nil)
			return
		}

		if room.Status == RoomStatusStart {
			for _, player := range room.Players {
				if playerID == player.ID {
					ctx.HTML(http.StatusOK, "playing.html", struct {
						PlayerID string
						Topic    string
						Room     *Room
					}{
						PlayerID: player.ID,
						Topic:    player.Topic,
						Room:     room,
					})
					return
				}
			}

			ctx.HTML(http.StatusOK, "start.html", struct {
				Room *Room
			}{
				Room: room,
			})
			return
		}

		if room.Status == RoomStatusPrepare {
			for _, player := range room.Players {
				if playerID == player.ID {
					ctx.HTML(http.StatusOK, "ready.html", struct {
						PlayerID string
						Topic    string
						Room     *Room
					}{
						PlayerID: player.ID,
						Topic:    player.Topic,
						Room:     room,
					})
					return
				}
			}

		}
	}

	ctx.HTML(http.StatusOK, "delete.html", nil)
}

func addPlayer(ctx *gin.Context) {
	roomID := ctx.Params.ByName("roomID")
	name := ctx.Params.ByName("name")

	if room := getRoom(roomID); room != nil {
		player := CreatePlayer(name)
		if err := room.addPlayer(player); err != nil {
			ctx.JSON(400, "人數已滿")
			return
		}

		ctx.JSON(Status_OK, struct {
			PlayerID string
			Topic    string
		}{
			PlayerID: player.ID,
			Topic:    player.Topic,
		})
		return
	}

	ctx.JSON(400, "不明原因壞了")
}

func getCard(ctx *gin.Context) {
	roomID := ctx.Params.ByName("roomID")
	playerID := ctx.Params.ByName("playerID")

	if room := getRoom(roomID); room != nil {
		for _, player := range room.Players {
			if playerID == player.ID {
				byteRoom, _ := json.Marshal(room)
				ctx.JSON(http.StatusOK, struct {
					PlayerID string
					Topic    string
					Dead     bool
					Room     string
				}{
					PlayerID: player.ID,
					Topic:    player.Topic,
					Dead:     player.Dead,
					Room:     string(byteRoom),
				})
				return
			}
		}
	}
}

func getRoom(roomID string) *Room {
	if val, ok := roomList.Load(roomID); ok {
		room := val.(*Room)
		return room
	}

	return nil
}

func gameStart(roomID string) {
	if room := getRoom(roomID); room != nil {
		room.start()
	}
}

func getTopic(index int) []string {
	return topic[index]
}

func getTopicIndex() int {
	return MapRandomKeyGet(topic).(int)
}

func MapRandomKeyGet(mapI interface{}) interface{} {
	keys := reflect.ValueOf(mapI).MapKeys()

	return keys[rand.Intn(len(keys))].Interface()
}
