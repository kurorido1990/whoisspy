package whoisspy

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

var gen Generator
var node *snowflake.Node
var roomList *sync.Map

func test(c *gin.Context) {
	data := new(IndexData)
	data.Title = "首頁"
	data.Content = "我的第一個首頁"
	c.HTML(http.StatusOK, "index.html", data)
}

func Run() {
	server := gin.Default()
	server.LoadHTMLGlob("template/*")
	server.GET("/", test)
	server.GET("/ping", func(c *gin.Context) {
		c.JSON(Status_OK, gin.H{
			"message": "pong",
		})
	})

	server.GET("/newRoom", newRoom)
	server.GET("/createRoom/:maxLimit", func(ctx *gin.Context) {
		maxLimit, _ := strconv.ParseInt(ctx.Params.ByName("maxLimit"), 10, 64)

		roomID := createRoom(int(maxLimit))

		ctx.JSON(Status_OK, gin.H{
			"addPlayer":   fmt.Sprintf("https://www.herokuapp.com/newPlayer/%d", roomID),
			"monitorRoom": fmt.Sprintf("https://www.herokuapp.com/room/%d", roomID),
		})
	})

	server.GET("/newPlayer/:roomID", newPlayer)
	server.GET("/addPlayer/:roomID/:name", addPlayer)

	server.GET("/room/:roomID/:playerID", gamePage)

	server.GET("/kick/:roomID/:playerID", kickPlayer)

	initSnowflake()
	initGen()
	server.Run(":9999")
	//server.Run()
}

func initSnowflake() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}

func initGen() {
	gen = CreateGen()
}

func newRoom(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "newRoom.html", nil)
}

func newPlayer(ctx *gin.Context) {
	roomID, _ := strconv.ParseInt(ctx.Params.ByName("roomID"), 10, 64)
	ctx.HTML(http.StatusOK, "newPlayer.html", struct {
		RoomID int64
	}{
		RoomID: roomID,
	})
}

func kickPlayer(ctx *gin.Context) {
	roomID, _ := strconv.ParseInt(ctx.Params.ByName("roomID"), 10, 64)
	playerID, _ := strconv.ParseInt(ctx.Params.ByName("playerID"), 10, 64)

	if room := getRoom(roomID); room != nil {
		room.kickPlayer(playerID)
	}

	ctx.JSON(Status_OK, nil)
}

func gamePage(ctx *gin.Context) {
	roomID, _ := strconv.ParseInt(ctx.Params.ByName("roomID"), 10, 64)
	playerID, _ := strconv.ParseInt(ctx.Params.ByName("playerID"), 10, 64)

	if room := getRoom(roomID); room != nil {
		if room.Status == RoomStatusEnd {
			ctx.HTML(http.StatusOK, "end.html", nil)
			return
		}

		if room.Status == RoomStatusStart {
			for _, player := range room.Players {
				if playerID == player.ID {
					ctx.HTML(http.StatusOK, "playing.html", struct {
						PlayerID int64
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
						PlayerID int64
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
	roomID, _ := strconv.ParseInt(ctx.Params.ByName("roomID"), 10, 64)
	name := ctx.Params.ByName("name")

	if room := getRoom(roomID); room != nil {
		player := CreatePlayer(name)
		if err := room.addPlayer(player); err != nil {
			ctx.JSON(400, err)
			return
		}

		ctx.JSON(Status_OK, struct {
			PlayerID int64  `json:"player_id"`
			Topic    string `json:"topic"`
		}{
			PlayerID: player.ID,
			Topic:    player.Topic,
		})
		return
	}

	ctx.HTML(http.StatusOK, "error.html", nil)
	return
}

func getRoom(roomID int64) *Room {
	if val, ok := roomList.Load(roomID); ok {
		room := val.(*Room)
		return room
	}

	return nil
}

func gameStart(roomID int64) {
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
