chain-indexer/
├── cmd/
│   ├── indexer/
│   │   └── main.go                # 程序入口（-p 参数、启动顺序）
│   ├── migrator/
│   │   └── main.go                # DB 初始化 / 迁移
│   └── tools/
│       └── shard-inspect.go       # shard 调试工具（可选）
│
├── internal/
│   ├── config/
│   │   └── config.go              # 全局配置、路径解析
│   │
│   ├── project/
│   │   ├── model.go               # Project / Checkpoint 数据模型
│   │   ├── repository.go          # projects / checkpoint 表访问
│   │   └── validator.go           # project 参数校验
│   │
│   ├── runtime/
│   │   ├── manager.go             # ProjectManager（唯一控制面）
│   │   ├── runtime.go             # ProjectRuntime（单项目运行态）
│   │   ├── control.go             # ControlMessage / ControlType
│   │   └── status.go              # RuntimeStatus 定义
│   │
│   ├── scanner/
│   │   ├── worker.go              # ScanWorker 主循环
│   │   ├── fetcher.go             # RPC 区块 / 日志获取
│   │   ├── block_cache.go         # block_number -> block_time 缓存
│   │   └── filter.go              # address / topics 构造
│   │
│   ├── storage/
│   │   ├── shard/
│   │   │   ├── router.go          # 时间 → shard 路由
│   │   │   ├── manager.go         # shard 生命周期管理
│   │   │   └── layout.go          # YYYY-MM 目录规则
│   │   │
│   │   ├── transaction/
│   │   │   ├── writer.go          # tx 写入接口
│   │   │   ├── reader.go          # tx 读取接口（查询层用）
│   │   │   └── model.go           # transaction 表结构
│   │   │
│   │   └── log/
│   │       ├── writer.go          # transaction log 写入
│   │       └── model.go
│   │
│   ├── query/
│   │   ├── engine.go              # Query Engine（HTTP / gRPC 共用）
│   │   ├── planner.go             # 时间范围 → shard 扫描计划
│   │   ├── cursor.go              # Cursor 编解码
│   │   └── paginator.go           # 分页控制
│   │
│   ├── api/
│   │   ├── http/
│   │   │   ├── server.go          # HTTP Server 启动
│   │   │   ├── middleware.go
│   │   │   ├── project_handler.go # add/update/active/status
│   │   │   └── query_handler.go   # tx 查询（HTTP）
│   │   │
│   │   └── grpc/
│   │       ├── server.go          # gRPC Server
│   │       ├── service.go         # TxQueryService 实现
│   │       └── interceptor.go
│   │
│   ├── proto/
│   │   └── tx_query.proto         # gRPC 协议定义
│   │
│   └── common/
│       ├── errors.go
│       ├── timeutil.go
│       ├── encoding.go
│       └── lifecycle.go
│
├── data/
│   └── {projectID}/
│       ├── meta.db                # projects / checkpoint
│       └── 2025-01/
│           ├── shard_000.db
│           └── shard_001.db
│
├── scripts/
│   ├── init_db.sh
│   └── cleanup.sh
│
├── go.mod
└── README.md
