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
	const future = 4070880000
	// 注册时间
	var regSince = 1512841140
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
		name           string
		timeZone       int
		openCount      int
		killCount      int
		RaidedSummary  string
		RaidedFrom     string
		StealedSummary string
		area           int
		rushCount      int
		score          string
	}

	var userInfoMap = make(map[int]userInfo)

	var totalUsers = 0
	// var registerSceneCount, upgradePartSceneCount = 0.0, 0.0
	// var noEnergySceneCount, inviteFriendSceneCount = 0.0, 0.0
	// var sentInvitationUsersCount, sentInvitationCount = 0.0, 0.0
	// var inviteSuccessUserCount, inviteSuccessCount = 0.0, 0.0
	// var area2UsersCount, area2facebookUsersCount = 0.0, 0.0
	//逐个用户处理
	for i := range users {
		var user = gj.GetIndex(i)
		var id = user.Get("userId").MustInt()
		var name = user.Get("user").Get("name").MustString()
		var timeZone = user.Get("user").Get("timeZone").MustInt()
		var registerTime = user.Get("user").Get("registerTime").MustInt()

		fmt.Println(registerTime)

		if registerTime < regSince || registerTime > regEnd { // 排除注册时间范围以外的用户，以时间戳为划分依据
			continue
		}

		if timeZone < -10 || timeZone > -4 { // 跳过东8区用户，不分析
			continue
		}
		totalUsers++

		var sequenceId = user.Get("area").Get("mode").Get("sequenceId").MustInt()
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

		//fmt.Println(id, name, timeZone)
		userInfoMap[id] = userInfo{
			name, timeZone, 0, 0, "", "", "", area, 0, "",
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
	var area2Rush3UsersCount, area2Rush4UsersCount, area2Rush5UsersCount, area2RushOtherUsersCount = 0, 0, 0, 0
	for k, v := range userInfoMap {
		_ = k
		fmt.Printf("%d\t%d\t%d\t%d\t%s\t%s\t%s\n", k, v.timeZone, v.openCount, v.killCount, v.RaidedSummary, v.StealedSummary, v.RaidedFrom)
		switch v.area {
		case 1:
			area1UsersCount++
			switch v.rushCount {
			case 0:
				area1Rush0UsersCount++
			case 1:
				area1Rush1UsersCount++
			case 2:
				area1Rush2UsersCount++
			case 3:
				area1Rush3UsersCount++
			case 4:
				area1Rush4UsersCount++
			default:
				area1RushOtherUsersCount++
			}
		case 2:
			area2UsersCount++
			switch v.rushCount {
			case 3:
				area2Rush3UsersCount++
			case 4:
				area2Rush4UsersCount++
				//fmt.Println(v.score)
			case 5:
				area2Rush5UsersCount++
			default:
				area2RushOtherUsersCount++
			}
		case 3:
			area3UsersCount++
		default:
			otherAreaUsersCount++
		}
	}
	fmt.Printf("新进用户数：%d\n", totalUsers)
	fmt.Printf("区域1流失总计：%d ，按飞船冲刺次数分布：\n", area1UsersCount)
	fmt.Printf("0次:%d\n1次:%d\n2次:%d\n3次:%d\n4次:%d\n其他:%d\n", area1Rush0UsersCount, area1Rush1UsersCount, area1Rush2UsersCount, area1Rush3UsersCount, area1Rush4UsersCount, area1RushOtherUsersCount)
	fmt.Printf("区域2流失总计：%d ，按飞船冲刺次数分布：\n", area2UsersCount)
	fmt.Printf("3次:%d\n4次:%d\n5次:%d\n其他:%d\n", area2Rush3UsersCount, area2Rush4UsersCount, area2Rush5UsersCount, area2RushOtherUsersCount)
	fmt.Printf("区域3流失总计：%d\n", area3UsersCount)
	fmt.Printf("其他区域流失总计：%d\n", otherAreaUsersCount)
}
