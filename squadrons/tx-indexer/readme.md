transaction-service/
├── cmd/
│   └── main.go                 # 启动服务
├── config/
│   └── config.yaml             # RPC 配置、项目配置
├── internal/
│   ├── scheduler/              # Scheduler Worker Pool
│   ├── storage/                # DAO + 分片路由
│   ├── api/                    # HTTP / gRPC API
│   ├── model/                  # Transaction / Log 数据模型
│   └── util/                   # DB utils, hash utils
└── pkg/
    └── project/                # 项目管理

+--------------------------------------------------------+
|                Transaction Storage Service            |
|                                                        |
|  +------------------+   +-------------------+         |
|  | Project A         |   | Project B         | ...     |
|  | Scheduler/Fetcher |   | Scheduler/Fetcher |         |
|  +--------+---------+   +---------+---------+         |
|           |                       |                   |
|           v                       v                   |
|  +-------------------------------------------+        |
|  |            Storage Module (DB)            |        |
|  |-------------------------------------------|        |
|  | transactions table / logs table           |        |
|  | 分片存储 / schema per project             |        |
|  +-----------------+-------------------------+        |
|                    |                                  |
|                    v                                  |
|  +-------------------------------------------+        |
|  |         Retrieval Module (API)            |        |
|  |-------------------------------------------|        |
|  | HTTP / gRPC 接口                           |        |
|  | 支持按 tx_hash, address, contract, time 查询|       |
|  +-------------------------------------------+        |
+--------------------------------------------------------+
