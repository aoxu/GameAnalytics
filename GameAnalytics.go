package main

import (
	"fmt"
	"io/ioutil"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

func main() {
	// 注册时间
	var regSince = 1512988980
	var regEnd = 1513189560

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
	fmt.Println("game data 表总样本数:", len(users))
	fmt.Println("statistic 表总样本数:", len(activity))

	type userInfo struct {
		name              string
		timeZone          int
		inviteFriendScene string
		NoEnergyScene     bool
		PermanentScene    bool
		UpgradePartScene  bool
	}

	var userInfoMap = make(map[int]userInfo)

	var totalUsers, facebookUsers = 0.0, 0.0
	var registerSceneCount, upgradePartSceneCount = 0.0, 0.0
	var noEnergySceneCount, inviteFriendSceneCount = 0.0, 0.0
	var sentInvitationUsersCount, sentInvitationCount = 0.0, 0.0
	var inviteSuccessUserCount, inviteSuccessCount = 0.0, 0.0
	var area2UsersCount, area2facebookUsersCount = 0.0, 0.0
	var likedUsersCount = 0.0
	//逐个用户处理
	for i, _ := range users {
		var user = gj.GetIndex(i)
		var id = user.Get("userId").MustInt()
		var name = user.Get("user").Get("name").MustString()
		var timeZone = user.Get("user").Get("timeZone").MustInt()
		var facebookId, err = user.Get("user").Get("facebookId").String()
		_ = err
		var registerTime = user.Get("user").Get("registerTime").MustInt()

		if registerTime < regSince || registerTime > regEnd { // 排除统计周期以外的用户，以时间戳为划分依据
			continue
		}

		if timeZone == 8 { // 跳过东8区用户，不分析
			continue
		}

		userInfoMap[id] = userInfo{
			name, timeZone, "", false, false, false,
		}

		//fmt.Println(id, name, timeZone)
		_ = id
		_ = name
		_ = timeZone
		totalUsers += 1

		var liked = user.Get("user").Get("facebookLikeReward").MustInt()
		if liked > 1 {
			likedUsersCount++
		}

		if facebookId == "" { // 不对未绑定用户进一步分析
			continue
		}
		//fmt.Println("facebookId", facebookId)
		facebookUsers += 1

		var bindScene = user.Get("user").Get("bindScene").MustString()
		var bindTime = user.Get("user").Get("bindTime").MustInt64()
		//fmt.Println(bindScene, bindTime)
		switch bindScene {
		case "register":
			registerSceneCount += 1
		case "upgradePart-Ship2-Part1-Level6":
			upgradePartSceneCount += 1
		case "upgradePart-Ship2-Part3-Level9":
			upgradePartSceneCount += 1
		case "upgradePart-Ship2-Part4-Level9":
			upgradePartSceneCount += 1
		case "upgradePart-Ship2-Part5-Level9":
			upgradePartSceneCount += 1
		case "freeEnergy":
			noEnergySceneCount += 1
		case "inviteFriend":
			inviteFriendSceneCount += 1
		default:
			fmt.Println("未知绑定场景：", bindScene)
		}

		var friend = user.Get("friend")
		var sequenceId = user.Get("area").Get("mode").Get("sequenceId").MustInt()

		pr, err := friend.Get("platformRequest").Array()
		var InviteRequestsCount = len(pr)
		if InviteRequestsCount > 0 {
			//fmt.Println("userId", id, "发出Facebook邀请", InviteRequestsCount, "份", "sequenceId=", sequenceId)
			sentInvitationUsersCount += 1.0
			sentInvitationCount = sentInvitationCount + float64(InviteRequestsCount)
		}

		if sequenceId >= 100 {
			area2UsersCount += 1
			if facebookId != "" {
				area2facebookUsersCount += 1
				strBindTime := time.Unix(bindTime, 0).Format("2006-01-02 15:04:05")
				//fmt.Println("player in area2", id, bindScene, strBindTime)
				_ = strBindTime
			}
		}

		var successInvitedCount = 0
		for i, _ := range pr {
			var inviteRequest = friend.Get("platformRequest").GetIndex(i)
			var complete = inviteRequest.Get("complete").MustBool()
			var invitedId = inviteRequest.Get("platformId").MustString()
			_ = invitedId
			if complete {
				successInvitedCount += 1
				//fmt.Println("发起邀请者id", id, "被邀请成功的人id是", invitedId)
			}
		}
		if successInvitedCount > 0 {
			inviteSuccessUserCount += 1.0
			inviteSuccessCount = inviteSuccessCount + float64(successInvitedCount)
		}
	}

	var NoEnergySceneCount = 0.0
	var UpgradePartSceneCount = 0.0
	var PermanentSceneCount = 0.0
	for i, _ := range activity {
		var act = dj.GetIndex(i)
		var userId = act.Get("userId").MustInt()
		//var time = act.Get("time").MustInt64()

		tmp, isExist := userInfoMap[userId]
		if !isExist { // user找不到，说明不在统计周期内，不统计
			continue
		}

		inviteFriendScene := act.Get("data").Get("inviteFriendScene")
		inviteFriendSceneArray, err := inviteFriendScene.Array()
		_ = err
		for j, _ := range inviteFriendSceneArray {
			var inviteScene = inviteFriendScene.GetIndex(j).Get("inviteScene").MustString()
			switch inviteScene {
			case "emptyEnergy":
				tmp.NoEnergyScene = true
			case "freeEnergy":
				tmp.NoEnergyScene = true
			case "upgradePart-Ship2-Part1-Level6":
				tmp.UpgradePartScene = true
			case "upgradePart-Ship2-Part3-Level9":
				tmp.UpgradePartScene = true
			case "upgradePart-Ship2-Part4-Level9":
				tmp.UpgradePartScene = true
			case "upgradePart-Ship2-Part5-Level9":
				tmp.UpgradePartScene = true
			case "upgradePart-Ship3-Part4-Level9":
				tmp.UpgradePartScene = true
			case "upgradePart-Ship3-Part5-Level10":
				tmp.UpgradePartScene = true
			case "upgradePart-Ship3-Part4-Level10":
				tmp.UpgradePartScene = true
			case "inviteForReward":
				tmp.PermanentScene = true
			case "inviteFriend":
				tmp.PermanentScene = true
			default:
				fmt.Println("未知邀请场景：", inviteScene)
			}
			tmp.inviteFriendScene += inviteScene + ","
		}

		userInfoMap[userId] = tmp
	}

	for k, v := range userInfoMap {
		if v.inviteFriendScene == "" {
			continue
		}
		//fmt.Printf("%d\t%d\t%s\n", k, v.timeZone, v.inviteFriendScene)
		_ = k
		if v.NoEnergyScene {
			NoEnergySceneCount += 1
		}
		if v.UpgradePartScene {
			UpgradePartSceneCount += 1
		}
		if v.PermanentScene {
			PermanentSceneCount += 1
		}
	}

	fmt.Println("新注册用户数: ", totalUsers)
	fmt.Println("绑定 Facebook 用户数: ", facebookUsers)
	fmt.Printf("Facebook 绑定率：%.2f%%\n", facebookUsers/totalUsers*100)
	fmt.Printf("游戏开始 场景绑定用户数：%.0f 占比：%.2f%%\n", registerSceneCount, registerSceneCount/facebookUsers*100)
	fmt.Printf("能量耗尽 场景绑定用户数：%.0f 占比：%.2f%%\n", noEnergySceneCount, noEnergySceneCount/facebookUsers*100)
	fmt.Printf("部件升级 场景绑定用户数：%.0f 占比：%.2f%%\n", upgradePartSceneCount, upgradePartSceneCount/facebookUsers*100)
	fmt.Printf("常驻入口 场景绑定用户数：%.0f 占比：%.2f%%\n", inviteFriendSceneCount, inviteFriendSceneCount/facebookUsers*100)
	fmt.Printf("发过邀请的用户数：%.0f 占所有新注册用户的比例：%.2f%%, 邀请了几个好友:%.0f\n", sentInvitationUsersCount, sentInvitationUsersCount/totalUsers*100, sentInvitationCount)
	fmt.Printf("能量耗尽 场景点过发送邀请按钮用户数：%.0f 占比：%.2f%%\n", NoEnergySceneCount, NoEnergySceneCount/sentInvitationUsersCount*100)
	fmt.Printf("部件升级 场景点过发送邀请按钮用户数：%.0f 占比：%.2f%%\n", UpgradePartSceneCount, UpgradePartSceneCount/sentInvitationUsersCount*100)
	fmt.Printf("常驻入口 场景点过发送邀请按钮用户数：%.0f 占比：%.2f%%\n", PermanentSceneCount, PermanentSceneCount/sentInvitationUsersCount*100)
	fmt.Printf("邀请成功的用户数：%.0f, 邀请成功了几个好友:%.0f\n", inviteSuccessUserCount, inviteSuccessCount)
	fmt.Printf("点赞用户数：%.0f, 占比:%.0f%%\n", likedUsersCount, likedUsersCount/totalUsers*100)
}
