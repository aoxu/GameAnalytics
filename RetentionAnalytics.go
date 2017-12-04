package main
 
import (
    "fmt"
    "io/ioutil"
    "strconv"
    simplejson "github.com/bitly/go-simplejson"
)
 
func main() {
    var timeSince int64 = 1512316800 // 跳过老用户，以时间戳为划分依据
    var timeEnd int64 = 1512403200

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
    // fmt.Println("game data 表总样本数:", len(users))
    // fmt.Println("statistic 表总样本数:", len(activity))

    type userInfo struct {
        name string
        timeZone int
        openCount int
        killCount int
        RaidedSummary string
        RaidedFrom string
        StealedSummary string
    }

    var userInfoMap = make(map[int]userInfo)

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
        var timeZone = user.Get("user").Get("timeZone").MustInt()
        var registerTime = user.Get("user").Get("registerTime").MustInt64()

        if registerTime < timeSince { // 跳过老用户，以时间戳为划分依据
            //continue;
        }

        if timeZone == 8 { // 跳过东8区用户，不分析
            continue;
        }

        //fmt.Println(id, name, timeZone)
        userInfoMap[id] = userInfo{
            name, timeZone, 0, 0, "", "", "",
        }
        //fmt.Println(id, userInfoMap[id].name, userInfoMap[id].timeZone)
    }

    for i, _ := range activity {
        var act = dj.GetIndex(i)
        var userId = act.Get("userId").MustInt()
        var time = act.Get("time").MustInt64()
        if time < timeSince || time > timeEnd  { // 排除统计周期外的用户数据
            continue;
        }

        tmp, isExist := userInfoMap[userId]
        if !isExist { // user找不到，说明是东八区用户，不统计
            continue;
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
        for j, _ := range revengeListArray {
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
        userInfoMap[userId] = tmp

        //fmt.Printf("%d\t%d\t%d\t%d\t%d\t%d\t%s\n", userId, tmp.timeZone, openCount, killCount, raidedCount, stealedCount, revengeIds)
    }

    for k, v := range userInfoMap {
        fmt.Printf("%d\t%d\t%d\t%d\t%s\t%s\t%s\n", k, v.timeZone, v.openCount, v.killCount, v.RaidedSummary, v.StealedSummary, v.RaidedFrom)
    }
}