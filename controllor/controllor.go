package controllor

import (
	"encoding/json"
	"fmt"
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

type BalanceRankingResponse struct {
	Ranking    int              `json:"ranking"`
	PlayerInfo model.PlayerInfo `json:"playerInfo"`
	Balance    float32          `json:"balance"`
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
	var res []BalanceRankingResponse
	for i, record := range ranking {
		res = append(res, BalanceRankingResponse{Ranking: i + 1 + (pageInt-1)*20, PlayerInfo: model.PlayerInfo{Uuid: record.UID, Name: record.Player, Avatar: con.service.GetAvatar(record.Player)}, Balance: record.Balance})
	}
	utils.SuccessResponse(c, "ok", res)
}

type PlayerKillResponse struct {
	Ranking              int    `json:"ranking"`
	PlayerId             string `json:"player_id"`
	PlayerName           string `json:"player_name"`
	AncientGuardianKills int    `json:"ancient_guardian_kills"`
	PhantomKills         int    `json:"phantom_kills"`
	PiglinBruteKills     int    `json:"piglin_brute_kills"`
	EnderDragonKills     int    `json:"ender_dragon_kills"`
	WitherKills          int    `json:"wither_kills"`
	WardenKills          int    `json:"warden_kills"`
	TotalKills           int    `json:"total_kills"`
	Avatar               string `json:"avatar"`
}

type PlayTimeResponse struct {
	Ranking    int    `json:"ranking"`
	UUID       string `json:"uuid"`
	PlayerName string `json:"player_name"`
	FirstLogin string `json:"first_login"`
	PlayTime   string `json:"play_time"`
	LastLogin  string `json:"last_login"`
	Avatar     string `json:"avatar"`
}
type FishAmountRankingResponse struct {
	Ranking    int    `json:"ranking"`
	PlayerName string `json:"player_name"`
	UUID       string `json:"uuid"`
	FishName   string `json:"fish_name"`
	Amount     int    `json:"amount"`
	Avatar     string `json:"avatar"`
}
type FishSizeRankingResponse struct {
	Ranking    int     `json:"ranking"`
	PlayerName string  `json:"player_name"`
	UUID       string  `json:"uuid"`
	FishName   string  `json:"fish_name"`
	Size       float32 `json:"size"`
	Avatar     string  `json:"avatar"`
}

func (con *Controller) getPlayerKillResponse(pageInt int, ranking []model.PlayerKillStats) []PlayerKillResponse {
	var res []PlayerKillResponse
	for i, record := range ranking {
		res = append(res, PlayerKillResponse{
			Ranking:              i + 1 + (pageInt-1)*20,
			PlayerId:             record.PlayerId,
			PlayerName:           con.service.GetNameByUUID(record.PlayerId),
			AncientGuardianKills: record.AncientGuardianKills,
			PhantomKills:         record.PhantomKills,
			PiglinBruteKills:     record.PiglinBruteKills,
			EnderDragonKills:     record.EnderDragonKills,
			WitherKills:          record.WitherKills,
			WardenKills:          record.WardenKills,
			TotalKills:           record.AncientGuardianKills + record.PhantomKills + record.PiglinBruteKills + record.EnderDragonKills + record.WitherKills + record.WardenKills,
			Avatar:               con.service.GetAvatar(con.service.GetNameByUUID(record.PlayerId)),
		})
	}
	return res
}
func formatTime(time int) string {
	hour := time / 3600
	minute := (time % 3600) / 60
	second := time % 60
	return strconv.Itoa(hour) + "小时 " + strconv.Itoa(minute) + "分钟 " + strconv.Itoa(second) + "秒"
}
func (con *Controller) HandleGetPlayerInfo(c *gin.Context) {
	username := c.Query("username")
	info, err := con.service.GetPlayerInfo(username)
	if err != nil {
		utils.ErrorResponse(c, 401, "error getting player info", err.Error())
		return
	}
	utils.SuccessResponse(c, "ok", info)
}
func (con *Controller) HandleGetPlayTime(c *gin.Context) {
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
	ranking := con.service.GetPlayTime(pageInt)
	var res []PlayTimeResponse
	for i, record := range ranking {
		res = append(res, PlayTimeResponse{
			Ranking:    i + 1 + (pageInt-1)*20,
			PlayerName: record.Name,
			UUID:       record.UUID,
			FirstLogin: record.FirstLogin.Format("2006-01-02 15:04:05"),
			PlayTime:   formatTime(record.PlayTime),
			LastLogin:  record.LastLogin.Format("2006-01-02 15:04:05"),
			Avatar:     con.service.GetAvatar(record.Name),
		})
	}
	utils.SuccessResponse(c, "ok", res)
}
func (con *Controller) HandleGetPlayerKillStatsSortByTotal(c *gin.Context) {
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
	ranking := con.service.GetPlayerKillStatsSortByTotal(pageInt)
	res := con.getPlayerKillResponse(pageInt, ranking)
	utils.SuccessResponse(c, "ok", res)
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
	res := con.getPlayerKillResponse(pageInt, ranking)
	utils.SuccessResponse(c, "ok", res)
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
	res := con.getPlayerKillResponse(pageInt, ranking)
	utils.SuccessResponse(c, "ok", res)
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
	res := con.getPlayerKillResponse(pageInt, ranking)
	utils.SuccessResponse(c, "ok", res)
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
	res := con.getPlayerKillResponse(pageInt, ranking)
	utils.SuccessResponse(c, "ok", res)
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
	res := con.getPlayerKillResponse(pageInt, ranking)
	utils.SuccessResponse(c, "ok", res)
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
	res := con.getPlayerKillResponse(pageInt, ranking)
	utils.SuccessResponse(c, "ok", res)
}
func (con *Controller) HandleGetFishRankingByAmount(c *gin.Context) {
	fish := c.Query("fish")
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return
	}
	records, err := con.service.GetFishRankingByAmount(fish, pageInt)
	if err != nil {
		utils.ErrorResponse(c, 401, "error getting fish data", "")
		return
	}
	var res []FishAmountRankingResponse
	for i, record := range records {
		res = append(res, FishAmountRankingResponse{Ranking: i + 1 + (pageInt-1)*20, PlayerName: con.service.GetNameByUUID(record.Uuid), UUID: record.Uuid, FishName: fish, Amount: record.Amount[fish], Avatar: con.service.GetAvatar(con.service.GetNameByUUID(record.Uuid))})
	}
	utils.SuccessResponse(c, "ok", res)
}
func (con *Controller) HandleGetFishRankingByTotalAmount(c *gin.Context) {
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return
	}
	fmt.Println("done1")
	records, err := con.service.GetFishRankingByTotalAmount(pageInt)
	if err != nil {
		utils.ErrorResponse(c, 401, "error getting fish data", "")
		return
	}
	fmt.Println("done2")
	var res []FishAmountRankingResponse
	for i, record := range records {
		res = append(res, FishAmountRankingResponse{Ranking: i + 1 + (pageInt-1)*20, PlayerName: con.service.GetNameByUUID(record.Uuid), UUID: record.Uuid, FishName: "", Amount: con.service.GetTotalAmount(record), Avatar: con.service.GetAvatar(con.service.GetNameByUUID(record.Uuid))})
	}
	fmt.Println("done3")
	utils.SuccessResponse(c, "ok", res)
}
func (con *Controller) HandleGetFishRankingBySize(c *gin.Context) {
	fish := c.Query("fish")
	page := c.Query("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		utils.ErrorResponse(c, 401, "invalid page number", "")
		return

	}
	records, err := con.service.GetFishRankingBySize(fish, pageInt)
	if err != nil {
		utils.ErrorResponse(c, 401, "error getting fish data", "")
		return
	}
	var res []FishSizeRankingResponse
	for i, record := range records {
		res = append(res, FishSizeRankingResponse{Ranking: i + 1 + (pageInt-1)*20, PlayerName: con.service.GetNameByUUID(record.Uuid), UUID: record.Uuid, FishName: fish, Size: record.MaxSize[fish], Avatar: con.service.GetAvatar(con.service.GetNameByUUID(record.Uuid))})
	}
	utils.SuccessResponse(c, "ok", res)
}
func (con *Controller) HandleGetPlayerList(c *gin.Context) {
	resp, err := http.Get("http://localhost:25577/api/players")
	if err != nil {
		utils.ErrorResponse(c, 501, "game server gg", "")
		return
	}

	defer resp.Body.Close()
	var response PlayerListResponse
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)

	if err != nil {
		utils.ErrorResponse(c, 401, "error getting player list", "")
		return
	}
	for i, player := range response.Lobby.Players {
		response.Lobby.Players[i].Avatar = con.service.GetAvatar(player.Name)
	}
	for i, player := range response.Survival.Players {
		response.Survival.Players[i].Avatar = con.service.GetAvatar(player.Name)
	}
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
	Ping   int    `json:"ping"`
	Name   string `json:"name"`
	UUID   string `json:"uuid"`
	Avatar string `json:"avatar"`
}

func (con *Controller) HandleGetPlayerProfileByName(c *gin.Context) {
	name := c.Query("name")
	profile := con.service.GetPlayerProfileByName(name)
	utils.SuccessResponse(c, "ok", profile)
}

func (con *Controller) HandleGetFish(c *gin.Context) {
	fishes, err := con.service.GetFish()
	if err != nil {
		utils.ErrorResponse(c, 401, "error getting fish data", err.Error())
		return
	}
	utils.SuccessResponse(c, "ok", fishes)
}
func (con *Controller) HandleGetLoginRecordCountByDate(c *gin.Context) {
	count := con.service.GetLoginRecordCountByDate()
	utils.SuccessResponse(c, "ok", count)
}
