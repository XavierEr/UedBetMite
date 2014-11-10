package main

import (
	"errors"
	"fmt"
	"github.com/XavierEr/UedBetMite/Model"
	"github.com/XavierEr/UedBetMite/UedBetDataJson"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var baseUri = "http://sb.uedbet.com/zh-cn/OddsService/"

func main() {
	oddsList, _ := scrapeOdds()
	insertOddsToMongoDb(oddsList)
}

func scrapeOdds() ([]model.Odds, error) {
	unixUtcDateTimeNowMillisecond := getUnixUtcDateTimeNowMillisecond()
	queryOddsParam := queryOddsParam{utcDateTime: unixUtcDateTimeNowMillisecond, sportId: 1, programmeId: 0, pageType: 1, uiBetType: "am", displayView: 2, pageNo: 0, oddsType: 2, sortBy: 1, isFirstLoad: true, MoreBetEvent: "null"}
	oddsUrl := getOddsUrl(queryOddsParam)

	uedBetData, err := getUedBetData(oddsUrl)
	if err != nil {
		return nil, err
	}

	oddsList := extractOdds("PreMatch", uedBetData.PreMatches)
	oddsList = append(oddsList, extractOdds("LiveMatch", uedBetData.LiveMatches)...)

	return oddsList, nil

	//~ fmt.Println(len(uedBetData.PreMatches.CategoryGroups))
	//~ Save odds data to mongo db here

	//~ for i := 1; i < uedBetData.TotalPages; i++ {
	//~ queryOddsParam.pageNo = i
	//~ oddsUrl = getOddsUrl(queryOddsParam)
	//~ uedBetData := getUedBetData(oddsUrl)
	//~ fmt.Println(uedBetData.TotalPages)
	//~ Save odds data to mongo db here
	//~ }
}

func extractOdds(matchType string, matches model.Match) []model.Odds {
	var oddsList []model.Odds

	for _, categoryGroup := range matches.CategoryGroups {
		for _, matchInfo := range categoryGroup.MatchInfos {
			homeRedCard, _ := strconv.Atoi(matchInfo.Information[8])
			awayRedCard, _ := strconv.Atoi(matchInfo.Information[9])
			homeScore, _ := strconv.Atoi(matchInfo.Information[10])
			awayScore, _ := strconv.Atoi(matchInfo.Information[11])

			oddsList = append(oddsList,
				model.Odds{
					OddsType:       matchType,
					LeagueKey:      categoryGroup.Category.Key,
					LeagueName:     categoryGroup.Category.Name,
					MatchKey:       matchInfo.Key,
					HomeTeamName:   matchInfo.Information[0],
					AwayTeamName:   matchInfo.Information[1],
					MatchStartDate: matchInfo.Information[4],
					MatchTime:      matchInfo.Information[5],
					HomeRedCard:    homeRedCard,
					AwayRedCard:    awayRedCard,
					HomeScore:      homeScore,
					AwayScore:      awayScore,
					MatchHalf:      matchInfo.Information[13]})
		}
	}
	return oddsList
}

func insertOddsToMongoDb(oddsList []model.Odds) {
	session, err := mgo.Dial("mongodb://localadmin:12qwer34@ds047040.mongolab.com:47040/uedbetmitedb")
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("uedbetmitedb").C("odds")

	for _, odds := range oddsList {
		err = c.Insert(odds)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getUedBetData(oddsUrl string) (model.UedBetData, error) {
	resp, err := http.Get(oddsUrl)
	if err != nil {
		return model.UedBetData{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return model.UedBetData{}, err
		}
		return uedBetDataJson.Parse(body), nil
	}
	return model.UedBetData{}, errors.New(string(resp.StatusCode))
}

func getOddsUrl(queryOddsParam queryOddsParam) string {
	return fmt.Sprintf("%vGetOdds?_=%v&sportId=%v&programmeId=%v&pageType=%v&uiBetType=%v&displayView=%v&pageNo=%v&oddsType=%v&sortBy=%v&isFirstLoad=%t&MoreBetEvent=%v",
		baseUri,
		queryOddsParam.utcDateTime,
		queryOddsParam.sportId,
		queryOddsParam.programmeId,
		queryOddsParam.pageType,
		queryOddsParam.uiBetType,
		queryOddsParam.displayView,
		queryOddsParam.pageNo,
		queryOddsParam.oddsType,
		queryOddsParam.sortBy,
		queryOddsParam.isFirstLoad,
		queryOddsParam.MoreBetEvent)
}

type queryOddsParam struct {
	utcDateTime                                                           int64
	sportId, programmeId, pageType, displayView, pageNo, oddsType, sortBy int
	uiBetType, MoreBetEvent                                               string
	isFirstLoad                                                           bool
}

func getUnixUtcDateTimeNowMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
