package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	simplejson "github.com/bitly/go-simplejson"
)

func main() {
	const releaseTime0_6_0 = 1510145520
	const releaseTime0_7_0 = 1510648680
	const releaseTime0_8_0 = 1511290800
	const releaseTime0_9_0 = 1512160080
	const releaseTime0_10_0 = 1512551940
	const releaseTime0_10_1 = 1512637740
	const releaseTime0_11_0 = 1512841140
	const bugfixed_0_11_0 = 1512988980
	const releaseTime0_12_0 = 1513189560
	const releaseTime0_13_0 = 1513362840
	const december16th = 1513353600
	const future = 4070880000
	// 注册时间
	var regSince = 1513353600
	var regEnd = 4070880000
	// 统计周期
	var statSince int64 = 0
	var statEnd int64 = 4070880000

	//gamedata.json
	gamedatajson, err := ioutil.ReadFile("gamedata.json")
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return
	}

	gj, err := simplejson.NewJson([]byte(gamedatajson))
	if err != nil {
		panic(err.Error())
	}

	users, err := gj.Array()

	//daily.json
	dailyjson, err := ioutil.ReadFile("statistic.json")
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return
	}

	dj, err := simplejson.NewJson([]byte(dailyjson))
	if err != nil {
		panic(err.Error())
	}

	activity, err := dj.Array()
	fmt.Println("gamedata 表总样本数:", len(users))
	fmt.Println("statistic 表总样本数:", len(activity))

	type userInfo struct {
		name              string
		timeZone          int
		openCount         int
		killCount         int
		RaidedSummary     string
		RaidedFrom        string
		StealedSummary    string
		area              int
		rushCount         int
		score             string
		device            string
		guide             string
		gold              int
		spaceshipLevels   string
		androidVersion    string
		dashCount         int
		explore           bool
		isClaimed         bool
		isExploreUnlocked bool
		isEnergyUsed      bool
		isLoginAgain      bool
	}

	var userInfoMap = make(map[int]userInfo)

	var totalUsers = 0
	var dailyUsers [20]int
	countryDailyUsers := make(map[string][]int)
	//逐个用户处理
	for i := range users {
		var user = gj.GetIndex(i)
		var timeZone = user.Get("user").Get("timeZone").MustInt()
		var registerTime = user.Get("user").Get("registerTime").MustInt()

		//fmt.Println(registerTime)

		if registerTime < regSince || registerTime > regEnd { // 排除注册时间范围以外的用户，以时间戳为划分依据
			continue
		}

		//if timeZone < 8 || timeZone > 10 {
		if (timeZone >= -10 && timeZone <= -2) || timeZone == 0 || timeZone == 1 || (timeZone >= 3 && timeZone <= 12) {
			//if timeZone != 8 {
			//continue
		}

		var dayIndex = (registerTime - regSince) / 86400
		dailyUsers[dayIndex]++
		totalUsers++

		var id = user.Get("userId").MustInt()
		var name = user.Get("user").Get("name").MustString()
		var device = user.Get("user").Get("phoneDevice").MustString()
		// var ip = user.Get("user").Get("ipInfo").Get("ip").MustString()
		var ipCountry = user.Get("user").Get("ipInfo").Get("ipcountry").MustString()
		var guide = user.Get("guide").MustMap()
		var gold = user.Get("resource").Get("gold").Get("count").MustInt()
		var spaceshipParts = user.Get("spaceShip").Get("spaceShips").Get("1").Get("parts")
		//spaceshipPartsMap, err := user.Get("spaceShip").Get("spaceShips").Get("1").Get("parts").Map()
		//_ = err
		var sequenceId = user.Get("area").Get("mode").Get("sequenceId").MustInt()
		var sequenceIndex = user.Get("area").Get("mode").Get("sequenceIndex").MustInt()
		var androidVersion = user.Get("user").Get("phoneSystemVer").MustString()
		var dashCount = user.Get("battle").Get("spaceShipDash").Get("count").MustInt()
		var isExploreUnlocked = user.Get("unlockSystem").Get("spaceShipExplore").MustBool()
		var lastLoginTime = user.Get("user").Get("lastLoginTime").MustInt()

		// if (country == "IND" && timeZone != 5) ||
		// 	(country == "CHN" && timeZone != 8) ||
		// 	(country == "THA" && timeZone != 7) ||
		// 	(country == "GBR" && timeZone != 0) ||
		// 	(country == "USA" && (timeZone < -10 || timeZone > -4)) {
		// 	fmt.Printf("%d\t%s\t%v\t%s\t%d\t%s\t%s\t%v\n", id, name, ipCountry, country, timeZone, device, androidVersion, ip)
		// }
		if ipCountry == "" {
			ipCountry = "未知"
		}
		if countryDailyUsers[ipCountry] == nil {
			countryDailyUsers[ipCountry] = append(countryDailyUsers[ipCountry], 0)
		}
		for dayIndex+1 > len(countryDailyUsers[ipCountry]) {
			countryDailyUsers[ipCountry] = append(countryDailyUsers[ipCountry], 0)
		}
		countryDailyUsers[ipCountry][dayIndex]++

		var area = 0
		if sequenceId < 100 {
			area = 1
		}
		if sequenceId > 100 && sequenceId < 200 {
			area = 2
		}
		if sequenceId > 200 && sequenceId < 300 {
			area = 3
		}
		if sequenceId > 300 {
			area = 4
		}
		var explore = false
		if sequenceIndex > 0 {
			explore = true
		}
		var isEnergyUsed = false
		if sequenceIndex >= 6 {
			isEnergyUsed = true
		}

		var isLoginAgain = false
		if lastLoginTime-registerTime > 3600 {
			isLoginAgain = true
		}

		var guideIDs = ""
		for k := range guide {
			guideIDs += k + ","
		}

		var spaceshipLevels = strconv.Itoa(spaceshipParts.Get("1").Get("level").MustInt()) + "," +
			strconv.Itoa(spaceshipParts.Get("2").Get("level").MustInt()) + "," +
			strconv.Itoa(spaceshipParts.Get("3").Get("level").MustInt()) + "," +
			strconv.Itoa(spaceshipParts.Get("4").Get("level").MustInt()) + "," +
			strconv.Itoa(spaceshipParts.Get("5").Get("level").MustInt())

		var isClaimed = spaceshipParts.Get("1").Get("claimed").MustBool()
		// for k := range spaceshipPartsMap {
		// 	spaceshipLevels += strconv.Itoa(spaceshipParts.Get(k).Get("level").MustInt()) + ","
		// }

		//fmt.Println(id, name, timeZone)
		userInfoMap[id] = userInfo{
			name, timeZone, 0, 0, "", "", "", area, 0, "", device, guideIDs, gold, spaceshipLevels, androidVersion, dashCount, explore, isClaimed, isExploreUnlocked, isEnergyUsed, isLoginAgain,
		}
		//fmt.Println(id, userInfoMap[id].name, userInfoMap[id].timeZone)
	}

	for i := range activity {
		var act = dj.GetIndex(i)
		var userId = act.Get("userId").MustInt()
		var time = act.Get("time").MustInt64()
		if time < statSince || time > statEnd { // 排除统计周期外的用户数据
			continue
		}

		tmp, isExist := userInfoMap[userId]
		if !isExist { // user找不到，说明是东八区用户，不统计
			continue
		}

		var friendBossTime = act.Get("data").Get("friendBossTime")
		openTime, err := friendBossTime.Get("openTime").Array()
		_ = err
		openCount := len(openTime)
		tmp.openCount += openCount

		killTime, err := friendBossTime.Get("killTime").Array()
		killCount := len(killTime)
		tmp.killCount += killCount

		var revengeIds = ""
		revengeList := act.Get("data").Get("attackedInfo").Get("revengeList")
		revengeListArray, err := revengeList.Array()
		for j := range revengeListArray {
			var char = ","
			if j == 0 {
				char = "/"
			} else {
				char = ","
			}
			raiderId, err := revengeList.GetIndex(j).Get("userId").Int()
			if err != nil {
				raiderId := revengeList.GetIndex(j).Get("userId").MustString()
				revengeIds += char + raiderId
				//fmt.Println("raider id(string): ", j, raiderId)
			} else {
				revengeIds += char + strconv.Itoa(raiderId)
				//fmt.Println("raider id(int): ", j, raiderId)
			}
		}
		tmp.RaidedFrom += revengeIds

		var raidedCount, stealedCount = 0, 0
		raidedCount += act.Get("data").Get("attackedInfo").Get("surpriseAttackedCount").MustInt()
		stealedCount += act.Get("data").Get("attackedInfo").Get("pvpAttackedCount").MustInt()
		tmp.RaidedSummary += strconv.Itoa(raidedCount) + "+"
		tmp.StealedSummary += strconv.Itoa(stealedCount) + "+"

		var rushScore = act.Get("data").Get("spaceShipDashScore")
		rushScoreArray, err := rushScore.Array()
		tmp.rushCount += len(rushScoreArray)
		for k := range rushScoreArray {
			tmp.score += rushScore.GetIndex(k).MustString() + "\t"
		}

		userInfoMap[userId] = tmp

		//fmt.Printf("%d\t%d\t%d\t%d\t%d\t%d\t%s\n", userId, tmp.timeZone, openCount, killCount, raidedCount, stealedCount, revengeIds)
	}

	var area1UsersCount, area2UsersCount, area3UsersCount, otherAreaUsersCount = 0, 0, 0, 0
	var area1Rush0UsersCount, area1Rush1UsersCount, area1Rush2UsersCount, area1Rush3UsersCount, area1Rush4UsersCount, area1RushOtherUsersCount = 0, 0, 0, 0, 0, 0
	var area2Rush0UsersCount, area2Rush1UsersCount, area2Rush2UsersCount, area2Rush3UsersCount, area2Rush4UsersCount, area2Rush5UsersCount, area2Rush6UsersCount, area2RushOtherUsersCount = 0, 0, 0, 0, 0, 0, 0, 0
	for k, v := range userInfoMap {
		_ = k
		//fmt.Printf("%d\t%d\t%d\t%d\t%s\t%s\t%s\n", k, v.timeZone, v.openCount, v.killCount, v.RaidedSummary, v.StealedSummary, v.RaidedFrom)
		switch v.area {
		case 1:
			area1UsersCount++
			switch v.dashCount {
			case 0:
				area1Rush0UsersCount++
			case 1:
				area1Rush1UsersCount++
				//fmt.Printf("剩余金币:\t%d\t引导:\t%s\t飞船部件等级:\t%s\t是否二次上线且间隔超过1小时:\t%v\n", v.gold, v.guide, v.spaceshipLevels, v.isLoginAgain)
			case 2:
				area1Rush2UsersCount++
				//fmt.Printf("剩余金币:\t%d\t引导:\t%s\t飞船部件等级:\t%s\t是否解锁探索:\t%v\t%v\n", v.gold, v.guide, v.spaceshipLevels, v.isExploreUnlocked, v.explore)
			case 3:
				area1Rush3UsersCount++
			case 4:
				area1Rush4UsersCount++
			default:
				area1RushOtherUsersCount++
			}
		case 2:
			area2UsersCount++
			switch v.dashCount {
			case 0:
				area2Rush0UsersCount++
			case 1:
				area2Rush1UsersCount++
			case 2:
				area2Rush2UsersCount++
			case 3:
				area2Rush3UsersCount++
				//fmt.Printf("剩余金币:\t%d\t引导:\t%s\t飞船部件等级:\t%s\t是否探索过:\t%t\t是否消耗完能量:\t%v\n", v.gold, v.guide, v.spaceshipLevels, v.explore, v.isEnergyUsed)
			case 4:
				area2Rush4UsersCount++
				//fmt.Printf("剩余金币:\t%d\t引导:\t%s\t飞船部件等级:\t%s\t是否探索过:\t%t\t是否领取满级奖励：\t%v\t是否消耗完能量\t%v\t冲刺得分:\t%s\n", v.gold, v.guide, v.spaceshipLevels, v.explore, v.isClaimed, v.explore, v.score)
			case 5:
				area2Rush5UsersCount++
			case 6:
				area2Rush6UsersCount++
			default:
				area2RushOtherUsersCount++
			}
		case 3:
			area3UsersCount++
		default:
			otherAreaUsersCount++
		}
	}
	fmt.Println("每日新进用户数：")
	for i := range dailyUsers {
		fmt.Printf("%d\t", dailyUsers[i])
	}
	fmt.Printf("\n")

	for k, v := range countryDailyUsers {
		fmt.Printf("%v\t", k)
		for i := range v {
			fmt.Printf("%v\t", v[i])
		}
		fmt.Printf("\n")
	}
	//fmt.Println(ipCountry, dayIndex, countryDailyUsers[ipCountry][dayIndex])

	fmt.Printf("\n新进用户数：\t%d\n", totalUsers)
	fmt.Printf("区域1流失总计：\t%d\t ，按飞船冲刺次数分布：\n", area1UsersCount)
	fmt.Printf("0次:\t%d\n1次:\t%d\n2次:\t%d\n3次:\t%d\n4次:\t%d\n其他:\t%d\n", area1Rush0UsersCount, area1Rush1UsersCount, area1Rush2UsersCount, area1Rush3UsersCount, area1Rush4UsersCount, area1RushOtherUsersCount)
	fmt.Printf("区域2流失总计：\t%d\t ，按飞船冲刺次数分布：\n", area2UsersCount)
	fmt.Printf("0次:\t%d\n1次:\t%d\n2次:\t%d\n3次:\t%d\n4次:\t%d\n5次:\t%d\n6次:\t%d\n其他:\t%d\n", area2Rush0UsersCount, area2Rush1UsersCount, area2Rush2UsersCount, area2Rush3UsersCount, area2Rush4UsersCount, area2Rush5UsersCount, area2Rush6UsersCount, area2RushOtherUsersCount)
	fmt.Printf("区域3流失总计：%d\n", area3UsersCount)
	fmt.Printf("其他区域流失总计：%d\n", otherAreaUsersCount)
}
