leslie-blog-server/
│
├── cmd/
│   ├── api/
│   │   └── main.go // API 服务入口
│   │
│   └── worker/
│       └── main.go // 工作进程入口
│
├── internal/
│   │
│   ├── config/
│   │   └── config.go // 配置文件
│   │
│   ├── server/
│   │   └── server.go // 服务入口
│   │
│   ├── middleware/
│   │   ├── cors.go // 跨域中间件
│   │   ├── jwt.go // JWT 中间件
│   │   ├── logger.go // 日志中间件
│   │   ├── recovery.go // 异常恢复中间件
│   │   └── permission.go // 权限中间件
│   │
│   ├── router/
│   │   └── router.go // 路由配置
│   │
│   ├── database/
│   │   ├── mysql.go // MySQL 数据库配置
│   │   └── redis.go // Redis 数据库配置
│   │
│   ├── logger/
│   │   └── logger.go // 日志配置文件
│   │
│   ├── response/
│   │   └── response.go // 响应配置文件
│   │
│   ├── errors/
│   │   └── errors.go // 错误配置文件
│   │
│   ├── pkg/
│   │   ├── jwt/ // JWT 相关配置
│   │   ├── password/ // 密码相关配置
│   │   ├── ulid/ // ULID 相关配置
│   │   ├── validator/ // 验证相关配置
│   │   └── pagination/ // 分页相关配置
│   │
│   └── modules/
│       │
│       ├── auth/
│       │   ├── handler/ // 认证模块路由处理函数
│       │   ├── service/ // 认证模块服务层
│       │   ├── repository/ // 认证模块数据访问层
│       │   ├── model/ // 认证模块模型层
│       │   ├── dto/ // 认证模块数据传输对象层
│       │   └── routes.go // 认证模块路由配置
│       │
│       ├── user/
│       │   ├── handler/
│       │   ├── service/
│       │   ├── repository/
│       │   ├── model/
│       │   ├── dto/
│       │   └── routes.go
│       │
│       ├── post/
│       │   ├── handler/
│       │   ├── service/
│       │   ├── repository/
│       │   ├── model/
│       │   ├── dto/
│       │   └── routes.go
│       │
│       ├── category/
│       │   ├── handler/
│       │   ├── service/
│       │   ├── repository/
│       │   ├── model/
│       │   ├── dto/
│       │   └── routes.go
│       │
│       ├── tag/
│       │   ├── handler/
│       │   ├── service/
│       │   ├── repository/
│       │   ├── model/
│       │   ├── dto/
│       │   └── routes.go
│       │
│       ├── comment/
│       │   ├── handler/
│       │   ├── service/
│       │   ├── repository/
│       │   ├── model/
│       │   ├── dto/
│       │   └── routes.go
│       │
│       ├── project/
│       │   ├── handler/
│       │   ├── service/
│       │   ├── repository/
│       │   ├── model/
│       │   ├── dto/
│       │   └── routes.go
│       │
│       └── dashboard/
│           ├── handler/
│           ├── service/
│           ├── repository/
│           ├── dto/
│           └── routes.go
│
├── migrations/
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   ├── 000002_create_posts.up.sql
│   └── ...
│
├── configs/
│   ├── config.yaml
│   ├── config.local.yaml
│   └── config.prod.yaml
│
├── docs/
│   └── swagger/
│
├── scripts/
│   ├── dev.sh
│   └── migrate.sh
│
├── deployments/
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── .env
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md