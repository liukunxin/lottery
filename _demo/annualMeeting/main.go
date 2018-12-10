/**
 * 年会抽奖程序
 * 增加了互斥锁，线程安全
 * 基础功能：
 * 1 /import 导入参与名单作为抽奖的用户
 * 2 /lucky 从名单中随机抽取用户
 * 测试方法：
 * curl http://localhost:8080/
 * curl --data "users=quincy,quincy2" http://localhost:8080/import
 * curl http://localhost:8080/lucky
 */

package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"math/rand"
	"strings"
	"time"
)

var userList []string
func newApp() *iris.Application  {
	app := iris.New()
	mvc.New(app.Party("/")).Handle(&lotteryController{})
	return app
}

type lotteryController struct {
	Ctx iris.Context
}

// 抽奖的控制器
func main()  {
	app := newApp()

	userList = make([]string, 0)

	app.Run(iris.Addr(":8000"))
}

// GET http://localhost:8080/
func (c *lotteryController) Get() string {
	count := len(userList)
	return fmt.Sprintf("当前总共参与抽奖的用户数: %d\n", count)
}

// POST http://localhost:8080/import
func (c *lotteryController) PostImport() string  {
	strUsers := c.Ctx.FormValue("users")
	users := strings.Split(strUsers, ",")
	count1 := len(userList)
	for _, u := range users{
		u = strings.TrimSpace(u)
		if len(u) > 0 {
			//导入用户
			userList = append(userList, u)
		}
	}
	count2 := len(userList)
	//注意这是一个数组，当再次抽奖时，userList包含之前抽奖的人数，而导入数据是往数组后添加
	// 所以成功导入人数，是当前数组数减去之前的数组数
	return fmt.Sprintf("当前总共参与抽奖的用户数: %d，成功导入用户数: %d\n", count2, (count2 - count1))
}

// GET http://localhost:8080/lucky
func (c *lotteryController) GetLucky() string {
	count := len(userList)
	if count > 1 {
		seed := time.Now().UnixNano()                                // rand内部运算的随机数
		index := rand.New(rand.NewSource(seed)).Int31n(int32(count)) // rand计算得到的随机数
		user := userList[index]                                      // 抽取到一个用户
		userList = append(userList[0:index], userList[index+1:]...)  // 移除这个用户
		return fmt.Sprintf("当前中奖用户: %s, 剩余用户数: %d\n", user, count-1)
	} else if count == 1 {
		user := userList[0]
		userList = userList[0:0]
		return fmt.Sprintf("当前中奖用户: %s, 剩余用户数: %d\n", user, count-1)
	} else {
		return fmt.Sprintf("已经没有参与用户，请先通过 /import 导入用户 \n")
	}

}
