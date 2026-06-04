# 出租房收租管理系统

## 技术栈
- 后端：Go 1.21+ + Gin 框架
- 数据库：MySQL 8.0
- 前端：Vue 3 + Vite + Element Plus
- 收据图片生成：html2canvas

## 项目结构
```
house/
├── backend/
│   ├── main.go
│   ├── config.yaml
│   ├── go.mod / go.sum
│   ├── config/        — 配置读取
│   ├── model/         — 数据模型
│   ├── repository/    — 数据库操作
│   ├── service/       — 业务逻辑
│   ├── handler/       — HTTP 处理器
│   ├── router/        — 路由注册
│   └── middleware/    — 中间件
├── frontend/
│   ├── src/
│   │   ├── views/
│   │   ├── api/
│   │   ├── router/
│   │   └── components/
│   ├── package.json
│   └── vite.config.js
└── README.md
```

## 数据库表

### rooms 房间信息表
```sql
CREATE TABLE rooms (
    id INT PRIMARY KEY AUTO_INCREMENT,
    room_no VARCHAR(10) NOT NULL UNIQUE COMMENT '房间号 如101、201',
    rent_price DECIMAL(10,2) NOT NULL COMMENT '月租金',
    elec_price DECIMAL(10,2) NOT NULL COMMENT '电费单价(元/度)',
    water_price DECIMAL(10,2) NOT NULL COMMENT '水费单价(元/吨)',
    is_rented TINYINT NOT NULL DEFAULT 1 COMMENT '1=已出租 0=未出租',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### bills 月度账单表
```sql
CREATE TABLE bills (
    id INT PRIMARY KEY AUTO_INCREMENT,
    room_no VARCHAR(10) NOT NULL COMMENT '房间号',
    year INT NOT NULL COMMENT '账单年份',
    month INT NOT NULL COMMENT '账单月份',
    elec_last DECIMAL(10,2) COMMENT '上月电表读数',
    elec_current DECIMAL(10,2) COMMENT '本月电表读数',
    elec_usage DECIMAL(10,2) COMMENT '用电量=current-last',
    elec_cost DECIMAL(10,2) COMMENT '电费=用量*单价',
    water_last DECIMAL(10,2) COMMENT '上月底数',
    water_current DECIMAL(10,2) COMMENT '本月底数',
    water_usage DECIMAL(10,2) COMMENT '用水量',
    water_cost DECIMAL(10,2) COMMENT '水费',
    rent_cost DECIMAL(10,2) COMMENT '月租金（快照保存，不随后续修改变化）',
    total_cost DECIMAL(10,2) COMMENT '总费用：房租+电费+水费',
    is_paid TINYINT DEFAULT 0 COMMENT '0=未收费 1=已收费',
    paid_at DATETIME NULL COMMENT '收费时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_room_year_month (room_no, year, month)
);
```

## 后端 API 清单
1. GET /api/rooms - 获取房间列表（支持出租状态筛选）
2. PUT /api/rooms/:room_no - 修改房间价格、出租状态
3. POST /api/bills/generate - 按年月生成当月所有房间账单（自动填充上月读数）
4. GET /api/bills/list - 获取账单列表（按年月、房间号、收费状态筛选）
5. PUT /api/bills/:id - 修改水电读数（未收费才可修改）
6. PUT /api/bills/:id/pay - 标记已收费
7. PUT /api/bills/:id/unpay - 取消收费
8. GET /api/bills/:id/receipt - 获取收据数据（用于前端生成图片）

## 后端要求
- 项目结构分层：handler、service、repository、model、config、router
- 自动初始化MySQL表结构（首次启动自动建表）
- 首次启动自动插入36间房间基础数据（101-112, 201-212, 301-312 三层每层12间）
- 每个房间默认：租金1200元、电费1.0元/度、水费5.0元/吨
- 配置文件 config.yaml 可外部修改MySQL连接信息
- 接口统一返回格式：{code, msg, data}
- 完善参数校验、错误处理
- 支持跨域（CORS）
- 默认端口 8080

## 前端要求
- 页面：
  1. 房间管理 - 价格设置、出租状态切换
  2. 月度账单 - 录入水电读数、标记收费
  3. 收据生成 - html2canvas 截图下载
  4. 历史账单查询
- 自动计算用量和费用
- 已收费账单置灰不可编辑
- 适配电脑端

## 启动方式
- 后端：`cd backend && go run main.go`
- 前端：`cd frontend && npm install && npm run dev`

## 代码规范
- 所有提示文字用中文
- 后端统一返回格式 {code, msg, data}
- 保持代码简洁，注释清晰
