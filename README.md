# 出租房收租管理系统 🏠

36间出租房月度收租账单管理系统，前后端分离架构。

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.18+ + Gin 框架 |
| 数据库 | MySQL 8.0 |
| 前端 | Vue 3 + Vite + Element Plus |
| 收据生成 | html2canvas |

## 项目结构

```
house/
├── backend/                # 后端 Go 项目
│   ├── main.go             # 入口：自动建表、初始化36间房、启动服务
│   ├── config.yaml         # MySQL 配置（可外部修改）
│   ├── config/             # 配置读取
│   ├── model/              # 数据模型（Room、Bill）
│   ├── repository/         # 数据库操作层
│   ├── service/            # 业务逻辑层（账单生成、收费等）
│   ├── handler/            # HTTP 处理器
│   ├── router/             # 路由注册
│   └── middleware/         # CORS + 统一响应格式
├── frontend/               # 前端 Vue3 项目
│   ├── src/
│   │   ├── views/          # 页面组件
│   │   ├── api/            # axios API 封装
│   │   └── router/         # 路由配置
│   └── index.html
└── README.md
```

## 快速启动

### 1. 配置 MySQL

修改 `backend/config.yaml`：

```yaml
mysql:
  address: 127.0.0.1
  port: 3306
  user: root
  password: ""      # 改成你的密码
  dbname: house_rent
```

> 需要提前创建 `house_rent` 数据库：`CREATE DATABASE house_rent;`

### 2. 启动后端

```bash
cd backend
go run main.go
```
服务默认运行在 `http://localhost:8080`，首次启动会自动创建表结构和36间房间。

### 3. 启动前端

```bash
cd frontend
npm run dev
```
默认地址 `http://localhost:5173`，已配置代理转发 `/api` 到后端。

## 功能说明

### 🏠 房间管理
- 36间房间（101-112, 201-212, 301-312）
- 自定义：月租金、电费单价、水费单价
- 开关出租状态，未出租房间不生成账单

### 📊 月度账单
- 按月生成所有已出租房间账单
- 自动读取上月水电读数，计算用量和费用
- 支持修改水电读数（已收费不可修改）
- 标记/取消收费

### 🧾 收据生成
- 一键生成美观收据样式
- 支持下载为图片，可截图发给租客

### 📅 历史账单
- 按年月查询历史账单
- 按收费状态筛选

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/rooms | 房间列表（?rented=0/1） |
| PUT | /api/rooms/:room_no | 修改房间配置 |
| POST | /api/bills/generate | 生成当月账单 |
| GET | /api/bills/list | 账单列表 |
| PUT | /api/bills/:id | 修改水电读数 |
| PUT | /api/bills/:id/pay | 标记已收费 |
| PUT | /api/bills/:id/unpay | 取消收费 |
| GET | /api/bills/:id/receipt | 收据数据 |

## 默认房间数据

| 楼层 | 房间号 | 租金(月) | 电费单价 | 水费单价 |
|------|--------|----------|----------|----------|
| 1楼 | 101-112 | 1200元 | 1.0元/度 | 5.0元/吨 |
| 2楼 | 201-212 | 1200元 | 1.0元/度 | 5.0元/吨 |
| 3楼 | 301-312 | 1200元 | 1.0元/度 | 5.0元/吨 |
