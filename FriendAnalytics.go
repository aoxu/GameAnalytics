package main
 
import (
    "fmt"
    "io/ioutil"
    simplejson "github.com/bitly/go-simplejson"
)
 
func main() {
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
    fmt.Println("game data 表总样本数:", len(users))

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
    fmt.Println("statistic 表总样本数:", len(activity))

    // var totalUsers, facebookUsers = 0.0, 0.0
    // var registerSceneCount, upgradePartSceneCount = 0.0, 0.0
    // var noEnergySceneCount, inviteFriendSceneCount = 0.0, 0.0
    // var sentInvitationUsersCount, sentInvitationCount = 0.0, 0.0
    // var inviteSuccessUserCount, inviteSuccessCount = 0.0, 0.0
    // var area2UsersCount, area2facebookUsersCount = 0.0, 0.0
    //逐个用户处理
    for i, _ := range users {
        var user = gj.GetIndex(i)
        var id = user.Get("userId").MustInt()
        var name = user.Get("user").Get("name").MustString()
        var timeZone = user.Get("user").Get("timeZone").MustFloat64()
        var registerTime = user.Get("user").Get("registerTime").MustFloat64()

        if registerTime < 1508731200 { // 跳过老用户，以时间戳为划分依据
            continue;
        }

        if timeZone == 8 { // 跳过东8区用户，不分析
            continue;
        }

        //fmt.Println(id, name, timeZone)
        _ = id
        _ = name
        _ = timeZone
        // totalUsers += 1

        // if facebookId == "" { // 不对未绑定用户进一步分析
        //     //continue;
        // }
        // //fmt.Println("facebookId", facebookId)
        // facebookUsers += 1

        // var customize = facebookUser.Get("customize")
        // var bindScene = customize.Get("bindScene").MustString()
        // var bindTime = customize.Get("bindTime").MustInt64()
        //fmt.Println(bindScene, bindTime)
        // switch bindScene {
        // case "register":
        //     registerSceneCount += 1
        // case "upgradePart-Ship2-Part1-Level6":
        //     upgradePartSceneCount += 1
        // case "upgradePart-Ship2-Part3-Level9":
        //     upgradePartSceneCount += 1
        // case "upgradePart-Ship2-Part4-Level9":
        //     upgradePartSceneCount += 1
        // case "freeEnergy":
        //     noEnergySceneCount += 1
        // case "inviteFriend":
        //     inviteFriendSceneCount += 1
        // }

        var friend = user.Get("friend")
        friends, err := friend.Get("friends").Map()
        _ = err
        var friendsCount = len(friends)

        pendingCount := friend.Get("pendingCount").MustInt()
    
        pendings, err := friend.Get("pendings").Map()
        var pendingsCount = len(pendings)
    
        friendRequests, err := friend.Get("requests").Map()
        var friendRequestsCount = len(friendRequests)

        var sequenceId = user.Get("area").Get("mode").Get("sequenceId").MustInt()
        var sequenceIndex = user.Get("area").Get("mode").Get("sequenceIndex").MustInt()

        fmt.Printf("%d\t%s\t%d\t%d\t%d\t%d\t%d\t%d\n", id, name, friendsCount, pendingsCount, friendRequestsCount, sequenceId, sequenceIndex, pendingCount)
    
        // pr, err := friend.Get("platformRequest").Array()
        // var InviteRequestsCount = len(pr)
        // if InviteRequestsCount > 0 {
        //     fmt.Println("userId", id, "发出Facebook邀请", InviteRequestsCount, "份", "sequenceId=", sequenceId)
        //     sentInvitationUsersCount += 1.0
        //     sentInvitationCount = sentInvitationCount + float64(InviteRequestsCount)
        // }

        // if sequenceId >= 100 {
        //     area2UsersCount += 1
        //     if facebookId != "" {
        //         area2facebookUsersCount += 1
        //         fmt.Println("player in area2", bindScene, bindTime)
        //     }
        // }
    
        // var successInvitedCount = 0
        // for i, _ := range pr {
        //     var inviteRequest = friend.Get("platformRequest").GetIndex(i)
        //     var complete = inviteRequest.Get("complete").MustBool()
        //     var invitedId = inviteRequest.Get("platformId").MustString()
        //     if complete {
        //         successInvitedCount += 1
        //         fmt.Println("发起邀请者id", id, "被邀请成功的人id是", invitedId)
        //     }
        // } 
        // if successInvitedCount > 0 {
        //     inviteSuccessUserCount += 1.0
        //     inviteSuccessCount = inviteSuccessCount + float64(successInvitedCount)
        // }
    }

    for i, _ := range activity {
        var act = dj.GetIndex(i)
        var userId = act.Get("userId").MustInt64()
        var time = act.Get("time").MustInt64()
        var friendBossTime = act.Get("data").Get("friendBossTime")
        openTime, err := friendBossTime.Get("openTime").Array()
        _ = err
        openCount := len(openTime)
        killTime, err := friendBossTime.Get("killTime").Array()
        killCount := len(killTime)
        fmt.Printf("%d\t%d\t%d\t%d\n", userId, time, openCount, killCount)
    }

    // fmt.Println("新注册用户数: ", totalUsers)
    // //fmt.Println("留存到区域2用户数：", area2UsersCount, "其中绑定用户数：", area2facebookUsersCount)
    // fmt.Println("绑定 Facebook 用户数: ", facebookUsers)
    // fmt.Printf("Facebook 绑定率：%.2f%%\n", facebookUsers/totalUsers*100)
    // fmt.Printf("游戏开始绑定数量：%.0f 占比：%.2f%%\n", registerSceneCount, registerSceneCount/facebookUsers*100)
    // fmt.Printf("能量耗尽绑定数量：%.0f 占比：%.2f%%\n", noEnergySceneCount, noEnergySceneCount/facebookUsers*100)
    // fmt.Printf("部件升级绑定数量：%.0f 占比：%.2f%%\n", upgradePartSceneCount, upgradePartSceneCount/facebookUsers*100)
    // fmt.Printf("常驻入口绑定数量：%.0f 占比：%.2f%%\n", inviteFriendSceneCount, inviteFriendSceneCount/facebookUsers*100)
    // fmt.Printf("发过邀请的用户数：%.0f, 邀请了几个好友:%.0f\n", sentInvitationUsersCount, sentInvitationCount)
    // fmt.Printf("邀请成功的用户数：%.0f, 邀请成功了几个好友:%.0f\n", inviteSuccessUserCount, inviteSuccessCount)
}