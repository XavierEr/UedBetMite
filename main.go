package main

import (
	"fmt"
	"github.com/XavierEr/UedBetMite/Model"
	"github.com/XavierEr/UedBetMite/UedBetDataJson"
	"io/ioutil"
	"net/http"
	"time"
)

var baseUri = "http://sb.uedbet.com/zh-cn/OddsService/"

func main() {
	unixUtcDateTimeNowMillisecond := getUnixUtcDateTimeNowMillisecond()
	queryOddsParam := queryOddsParam{utcDateTime: unixUtcDateTimeNowMillisecond, sportId: 1, programmeId: 0, pageType: 1, uiBetType: "am", displayView: 2, pageNo: 0, oddsType: 2, sortBy: 1, isFirstLoad: true, MoreBetEvent: "null"}
	oddsUrl := getOddsUrl(queryOddsParam)

	uedBetData := getUedBetData(oddsUrl)
	fmt.Println(uedBetData.TotalPages)
	//~ Save odds data to mongo db here

	for i := 1; i < uedBetData.TotalPages; i++ {
		queryOddsParam.pageNo = i
		oddsUrl = getOddsUrl(queryOddsParam)
		uedBetData := getUedBetData(oddsUrl)
		fmt.Println(uedBetData.TotalPages)
		//~ Save odds data to mongo db here
	}
}

func getUedBetData(oddsUrl string) model.UedBetData {
	resp, err := http.Get(oddsUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	//~ if resp.StatusCode == 200 {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return uedBetDataJson.Parse(body)
	//~ }
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
