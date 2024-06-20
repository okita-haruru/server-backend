package controllor

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"sushi/model"
	"sushi/service"
	"sushi/utils"
	"sushi/utils/config"
)

type Controller struct {
	service *service.Service
	log     *logrus.Logger
	conf    *config.Config
}

func NewControllor(service *service.Service, log *logrus.Logger, conf *config.Config) *Controller {
	return &Controller{service: service, log: log, conf: conf}
}

func (con *Controller) HandlePing(c *gin.Context) {

	con.log.Debug("handling ping...")
	utils.SuccessResponse(c, "pong", "")
}
func (con *Controller) HandleGetBalance(c *gin.Context) {
	uuid := c.Query("uuid")
	balance := con.service.GetBalance(uuid)
	utils.SuccessResponse(c, "ok", balance)
}
func (con *Controller) HandleGetBalanceRanking(c *gin.Context) {
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return

	}
	if pageInt == -1 {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return
	}
	ranking := con.service.GetBalanceRanking(pageInt)
	utils.SuccessResponse(c, "ok", ranking)
}
func (con *Controller) HandleGetPlayerKillStatsSortByWarden(c *gin.Context) {
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return

	}
	if pageInt == -1 {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return
	}
	ranking := con.service.GetPlayerKillStatsSortByWarden(pageInt)
	utils.SuccessResponse(c, "ok", ranking)
}
func (con *Controller) HandleGetPlayerKillStatsSortByEnderDragon(c *gin.Context) {
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return

	}
	if pageInt == -1 {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return
	}
	ranking := con.service.GetPlayerKillStatsSortByEnderDragon(pageInt)
	utils.SuccessResponse(c, "ok", ranking)
}
func (con *Controller) HandleGetPlayerKillStatsSortByWither(c *gin.Context) {
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return

	}
	if pageInt == -1 {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return
	}
	ranking := con.service.GetPlayerKillStatsSortByWither(pageInt)
	utils.SuccessResponse(c, "ok", ranking)
}
func (con *Controller) HandleGetPlayerKillStatsSortByPiglin(c *gin.Context) {
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return

	}
	if pageInt == -1 {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return
	}
	ranking := con.service.GetPlayerKillStatsSortByPiglin(pageInt)
	utils.SuccessResponse(c, "ok", ranking)
}
func (con *Controller) HandleGetPlayerKillStatsSortByPhantom(c *gin.Context) {
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return

	}
	if pageInt == -1 {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return
	}
	ranking := con.service.GetPlayerKillStatsSortByPhantom(pageInt)
	utils.SuccessResponse(c, "ok", ranking)
}
func (con *Controller) HandleGetPlayerKillStatsSortByAncientGuardian(c *gin.Context) {
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return

	}
	if pageInt == -1 {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return
	}
	ranking := con.service.GetPlayerKillStatsSortByAncientGuardian(pageInt)
	utils.SuccessResponse(c, "ok", ranking)
}
func (con *Controller) HandleGetFishAmount(c *gin.Context) {
	err, records := con.service.GetDecodedFishData()
	if err != nil {
		utils.ErrorResponse(c, 401, "error getting fish data", "")
		return
	}
	var fishAmountData []model.FishAmountData
	for _, record := range records {
		fishAmountData = append(fishAmountData, model.FishAmountData{Uuid: record.Uuid, Amount: record.Amount})
	}
	utils.SuccessResponse(c, "ok", fishAmountData)
}
func (con *Controller) HandleGetFishSize(c *gin.Context) {
	err, records := con.service.GetDecodedFishData()
	if err != nil {
		utils.ErrorResponse(c, 401, "error getting fish data", "")
		return
	}
	var fishSizeData []model.FishSizeData
	for _, record := range records {
		fishSizeData = append(fishSizeData, model.FishSizeData{Uuid: record.Uuid, MaxSize: record.MaxSize})
	}
	utils.SuccessResponse(c, "ok", fishSizeData)
}
func (con *Controller) HandleGetPlayerList(c *gin.Context) {
	resp, err := http.Get("http://localhost:25577/api/players")
	if err != nil {
		utils.ErrorResponse(c, 401, "error getting player list", "")
	}
	defer resp.Body.Close()
	var response PlayerListResponse
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	utils.SuccessResponse(c, "ok", response)
}

type PlayerListResponse struct {
	Lobby    RoomJson `json:"lobby"`
	Survival RoomJson `json:"survival"`
}
type RoomJson struct {
	Players []PlayerJson `json:"players"`
	Count   int          `json:"count"`
}
type PlayerJson struct {
	Ping int    `json:"ping"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

func (con *Controller) HandleGetPlayerProfileByName(c *gin.Context) {
	name := c.Query("name")
	profile := con.service.GetPlayerProfileByName(name)
	utils.SuccessResponse(c, "ok", profile)
}
