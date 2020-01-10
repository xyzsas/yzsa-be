package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"yzsa-be/controllers"
	"yzsa-be/middleware"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Default())

	api := router.Group("/api")
	{
		C := api.Group("/C")
		{
			// 用户登录
			C.GET("/auth", controllers.User.HandShake) // 登录第一次握手
			C.POST("/auth", controllers.User.Login)    // 登录
			C.DELETE("/auth", controllers.User.Logout) // 登出
		}

		U := api.Group("/U")
		U.Use(middleware.UserAuth())
		{
			// 用户操作
			U.PUT("/pwd", controllers.User.UpdatePwd)                                  // 密码修改
			U.GET("/task", controllers.Task.GetOpen)                                   // 用户对应任务列表
			U.GET("/task/:id", middleware.Task(false), controllers.Task.GetRealTime)   // 获取任务实时信息
			U.GET("/record/:id", middleware.Task(false), controllers.Record.GetOne)    // 获取任务个人记录
			U.POST("/record/:id", middleware.Task(false), controllers.Record.Response) // 响应任务
		}

		T := api.Group("/T")
		T.Use(middleware.UserAuth())
		T.Use(middleware.RoleAuth("teacher"))
		{
			T.GET("/task", controllers.Task.GetList)                                                    // 获取所有任务
			T.GET("/task/:id", middleware.Task(true), controllers.Task.GetOne)                          // 获取任务详细信息
			T.POST("/task", controllers.Task.Insert)                                                    // 新建任务
			T.PUT("/task/:id", middleware.Task(true), controllers.Task.Update)                          // 修改任务
			T.PUT("/task/:id/open", middleware.Task(true), controllers.Task.Open)                       // 开启任务
			T.PUT("/task/:id/close", middleware.Task(true), controllers.Task.Close)                     // 关闭任务
			T.DELETE("/task/:id", middleware.Task(true), controllers.Task.Delete)                       // 删除任务
			T.GET("/record/:id", middleware.Task(true), controllers.Record.GetAll)                      // 获取任务全部记录
			T.POST("/permission/:id", middleware.Task(true), controllers.Permission.AddTemp)            // 添加临时权限
			T.DELETE("/permission/:id/:user", middleware.Task(true), controllers.Permission.DeleteTemp) // 删除临时权限
		}

		A := api.Group("/A")
		A.Use(middleware.UserAuth())
		A.Use(middleware.RoleAuth("admin"))
		{
			// 对权限树的操作
			A.GET("/permission", controllers.Permission.GetList)       // 权限节点列表
			A.POST("/permission", controllers.Permission.Insert)       // 添加权限节点
			A.PUT("/permission", controllers.Permission.Update)        // 修改权限节点
			A.DELETE("/permission/:id", controllers.Permission.Delete) // 删除权限节点
			// 对用户的操作
			A.GET("/user", controllers.User.GetList)                         // 用户列表
			A.POST("/user", controllers.User.Insert)                         // 添加用户
			A.GET("/user/:id", controllers.User.GetOne)                      // 获取用户信息
			A.DELETE("/user/:id", controllers.User.Delete)                   // 删除用户
			A.PUT("/user/:id/permission", controllers.User.UpdatePermission) // 修改用户权限
			A.PUT("/user/:id/name", controllers.User.UpdateName)             // 修改用户名称
			A.PUT("/user/:id/pwd", controllers.User.ResetPwd)                // 重置用户密码
		}
	}

	_ = router.Run(":5000")
}
