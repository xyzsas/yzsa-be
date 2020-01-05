# API接口路由

> 所有API均在一级路由`/api`下

## Common

> 二级路由`/C`

所有用户可访问

| 相对路由 | HTTP方法 | 说明           |
| -------- | -------- | -------------- |
| `/auth`  | `GET`    | 登录第一次握手 |
|          | `POST`   | 登录第二次握手 |
|          | `DELETE` | 登出           |

## User

> 二级路由`/U`，中间件`UserAuth`

登录用户可访问

| 相对路由      | HTTP方法 | 说明             |
| ------------- | -------- | ---------------- |
| `/pwd`        | `PUT`    | 修改密码         |
| `/task`       | `GET`    | 获取任务列表     |
| `/task/:id`   | `GET`    | 获取任务实时信息 |
| `/record/:id` | `GET`    | 获取任务个人记录 |
|               | `POST`   | 响应任务         |

## Teacher

> 二级路由`/T`，中间件`UserAuth` `RoleAuth("teacher")`

教师可访问

| 相对路由                | HTTP方法 | 说明                 |
| ----------------------- | -------- | -------------------- |
| `/task`                 | `GET`    | 获取可管理任务列表   |
| `/task/:id`             | `GET`    | 获取任务详细信息     |
| `/task`                 | `POST`   | 新建任务             |
| `/task/:id`             | `PUT`    | 修改任务             |
| `/task/:id/open`        | `PUT`    | 开启任务             |
| `/task/:id/close`       | `PUT`    | 关闭任务             |
| `/task/:id`             | `DELETE` | 删除任务             |
| `/record/:id`           | `GET`    | 读取任务全部记录     |
| `/permission/:id`       | `POST`   | 添加临时权限         |
| `/permission/:id/:user` | `DELETE` | 删除临时权限         |

## Admin

> 二级路由`/A`，中间件`UserAuth` `RoleAuth("admin")`

管理员可访问

| 相对路由          | HTTP方法 | 说明           |
| ----------------- | -------- | -------------- |
| `/permission`     | `GET`    | 权限节点列表   |
|                   | `POST`   | 添加权限节点   |
|                   | `PUT`    | 修改权限节点   |
| `/permission/:id` | `DELETE` | 删除权限节点   |
| `/user`           | `GET`    | 用户列表       |
|                   | `POST`   | 添加用户       |
| `/admin/:id`      | `GET`    | 获取用户信息   |
|                   | `PUT`    | 修改用户权限   |
|                   | `DELETE` | 删除用户       |
| `/admin/:id/pwd`  | `PUT`    | 重置管理员密码 |