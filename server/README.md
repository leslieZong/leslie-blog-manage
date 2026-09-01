go-enterprise/
├── cmd
│   └── api                # 程序入口
│       └── main.go
├── internal               # 业务私有代码，外部无法导入
│   ├── api                # handler 控制器层
│   │   ├── handler.go
│   │   └── user.go
│   ├── service            # 业务逻辑层
│   │   └── user_service.go
│   ├── repository         # 数据访问层 DB操作
│   │   └── user_repo.go
│   ├── model              # 数据模型：数据库model、请求/响应DTO
│   │   ├── dto
│   │   │   └── user_dto.go
│   │   └── entity
│   │       └── user.go
│   └── middleware         # 内部中间件 jwt、auth、限流、日志
│       ├── jwt.go
│       ├── logger.go
│       └── cors.go
├── pkg                    # 公共可复用包，可以外部引用
│   ├── config             # viper配置解析
│   ├── db                 # mysql gorm初始化
│   ├── redis              # redis初始化
│   ├── logger             # zap日志封装
│   ├── jwt                # jwt工具
│   ├── response           # 统一返回封装
│   ├── validator          # 参数校验
│   └── utils              # 工具函数
├── configs                # 配置文件 yaml
│   └── app.yaml
├── scripts                # 脚本：sql、shell、docker
│   └── init.sql
├── api                    # openapi swagger文档
├── test                   # 单元测试
├── go.mod
├── go.sum
├── .gitignore
└── Dockerfile
