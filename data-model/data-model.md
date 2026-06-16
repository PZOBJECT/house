# 出租房收租管理系统 · 数据模型

## 概述

本系统包含两个核心数据实体：**房间（Room）** 和 **账单（Bill）**。采用 Go + GORM 作为 ORM 框架，MySQL 8.0 作为数据库。

---

## 一、房间模型（Room）

### 数据库表：`rooms`

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| `id` | INT | PRIMARY KEY, AUTO_INCREMENT | 自增主键 |
| `room_no` | VARCHAR(10) | UNIQUE, NOT NULL | 房间号，如 `101`、`201` |
| `rent_price` | DECIMAL(10,2) | NOT NULL | 月租金（元） |
| `elec_price` | DECIMAL(10,2) | NOT NULL | 电费单价（元/度） |
| `water_price` | DECIMAL(10,2) | NOT NULL | 水费单价（元/吨） |
| `is_rented` | TINYINT | NOT NULL, DEFAULT 1 | 出租状态：`1`=已出租，`0`=未出租 |
| `created_at` | DATETIME | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `updated_at` | DATETIME | ON UPDATE CURRENT_TIMESTAMP | 更新时间 |

### Go 结构体

```go
type Room struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    RoomNo     string    `gorm:"column:room_no;type:varchar(10);uniqueIndex;not null" json:"room_no"`
    RentPrice  float64   `gorm:"column:rent_price;type:decimal(10,2);not null" json:"rent_price"`
    ElecPrice  float64   `gorm:"column:elec_price;type:decimal(10,2);not null" json:"elec_price"`
    WaterPrice float64   `gorm:"column:water_price;type:decimal(10,2);not null" json:"water_price"`
    IsRented   int       `gorm:"column:is_rented;type:tinyint;not null;default:1" json:"is_rented"`
    CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
    UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}
```

### 初始数据

系统首次启动时自动插入 **36 间房间**，按 3 层楼分布：

| 楼层 | 房间号范围 | 间数 |
|------|-----------|------|
| 1 楼 | 101–112 | 12 |
| 2 楼 | 201–212 | 12 |
| 3 楼 | 301–312 | 12 |

默认值：租金 **1200 元**、电费 **1.0 元/度**、水费 **5.0 元/吨**，全部默认为已出租。

### 业务规则

- 房间号不可重复，格式为 `[楼层][编号]`，如 `101`（1 楼 01 号）
- `rent_price`、`elec_price`、`water_price` 可编辑，修改后**不影响已有账单**（账单快照保存价目）
- `is_rented=0` 的房间不会在账单生成中产生账单
- 楼层由 `room_no` 首字符推断（`1`=1 楼，`2`=2 楼，`3`=3 楼）

---

## 二、账单模型（Bill）

### 数据库表：`bills`

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| `id` | INT | PRIMARY KEY, AUTO_INCREMENT | 自增主键 |
| `room_no` | VARCHAR(10) | NOT NULL, UNIQUE INDEX | 房间号 |
| `year` | INT | NOT NULL, UNIQUE INDEX | 账单年份（如 `2026`） |
| `month` | INT | NOT NULL, UNIQUE INDEX | 账单月份（`1`–`12`） |
| `elec_last` | DECIMAL(10,2) | NULLABLE | 上月电表读数 |
| `elec_current` | DECIMAL(10,2) | NULLABLE | 本月电表读数 |
| `elec_usage` | DECIMAL(10,2) | NULLABLE | 用电量 = `elec_current - elec_last` |
| `elec_cost` | DECIMAL(10,2) | NULLABLE | 电费 = `elec_usage × elec_price` |
| `water_last` | DECIMAL(10,2) | NULLABLE | 上月水表底数 |
| `water_current` | DECIMAL(10,2) | NULLABLE | 本月水表底数 |
| `water_usage` | DECIMAL(10,2) | NULLABLE | 用水量 = `water_current - water_last` |
| `water_cost` | DECIMAL(10,2) | NULLABLE | 水费 = `water_usage × water_price` |
| `rent_cost` | DECIMAL(10,2) | NOT NULL | 月租金（快照，不随后续修改变化） |
| `total_cost` | DECIMAL(10,2) | NOT NULL | 总费用 = `rent_cost + elec_cost + water_cost` |
| `is_paid` | TINYINT | DEFAULT `0` | 收费状态：`0`=未收费，`1`=已收费 |
| `paid_at` | DATETIME | NULL | 收费时间（标记收费时自动填入） |
| `created_at` | DATETIME | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| `updated_at` | DATETIME | ON UPDATE CURRENT_TIMESTAMP | 更新时间 |

**唯一索引**：`uk_room_year_month`（`room_no + year + month`），确保每个房间每月只有一条账单。

### Go 结构体

```go
type Bill struct {
    ID           uint       `gorm:"primaryKey" json:"id"`
    RoomNo       string     `gorm:"column:room_no;..." json:"room_no"`
    Year         int        `gorm:"column:year;..." json:"year"`
    Month        int        `gorm:"column:month;..." json:"month"`
    ElecLast     *float64   `gorm:"column:elec_last;...” json:"elec_last"`
    ElecCurrent  *float64   `gorm:"column:elec_current;...” json:"elec_current"`
    ElecUsage    *float64   `gorm:"column:elec_usage;...” json:"elec_usage"`
    ElecCost     *float64   `gorm:"column:elec_cost;...” json:"elec_cost"`
    WaterLast    *float64   `gorm:"column:water_last;...” json:"water_last"`
    WaterCurrent *float64   `gorm:"column:water_current;...” json:"water_current"`
    WaterUsage   *float64   `gorm:"column:water_usage;...” json:"water_usage"`
    WaterCost    *float64   `gorm:"column:water_cost;...” json:"water_cost"`
    RentCost     float64    `gorm:"column:rent_cost;...” json:"rent_cost"`
    TotalCost    float64    `gorm:"column:total_cost;...” json:"total_cost"`
    IsPaid       int        `gorm:"column:is_paid;...” json:"is_paid"`
    PaidAt       *time.Time `gorm:"column:paid_at" json:"paid_at"`
    CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
    UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}
```

> 注：电费/水费相关的字段使用 `*float64`（指针类型），因为新增账单时这些字段为 `NULL`，需要手动录入读数后才会有值。

---

## 三、数据关系图

```
┌─────────────┐        ┌──────────────────────────────┐
│    Room     │        │           Bill               │
├─────────────┤        ├──────────────────────────────┤
│ room_no (PK)│──1:N──▶│ room_no                      │
│ rent_price  │        │ year + month (复合唯一)      │
│ elec_price  │        │ rent_cost (快照自 rent_price)│
│ water_price │        │ elec_cost (计算自 elec_usage)│
│ is_rented   │        │ water_cost (计算自 water_usage)│
└─────────────┘        │ total_cost (汇总)            │
                       │ is_paid / paid_at            │
                       └──────────────────────────────┘
```

**核心关系**：
- 一个房间每月产生 **一条** 账单记录（`1:N` 关系）
- 账单生成时**快照**房间的当前租房，后续房间调价不影响已有账单
- 账单的费用计算依赖房间的**单价**（仅生成时读取）

---

## 四、账单生命周期

```
┌──────────┐     ┌───────────┐     ┌──────────┐     ┌──────────┐
│ 生成账单  │────▶│ 录入读数   │────▶│ 标记收费  │────▶│ 取消收费  │
│ (自动)   │     │ (手动)    │     │ (手动)   │     │ (手动)   │
└──────────┘     └───────────┘     └──────────┘     └──────────┘
```

### 4.1 生成账单（POST `/api/bills/generate`）

- 输入：`year` + `month`
- 逻辑：
  1. 查询所有 `is_rented=1` 的房间
  2. 对每个房间，检查是否已有该年月的账单（唯一索引防止重复）
  3. 若无，则新建一条账单：
     - `elec_last` = 该房间**上月的 `elec_current`**（若无上月记录则置空）
     - `water_last` = 该房间**上月的 `water_current`**
     - `elec_current` / `water_current` 置空（待录入）
     - `elec_usage` / `water_usage` / `elec_cost` / `water_cost` 置空
     - `rent_cost` = 房间当前 `rent_price`（快照）
     - `total_cost` = `rent_cost`（水电费未录入时为纯房租）

### 4.2 录入读数（PUT `/api/bills/:id`）

- 输入：`elec_current` + `water_current`
- 约束：仅 `is_paid=0`（未收费）的账单可修改
- 计算：
  - `elec_usage = elec_current - elec_last`（从房间表读取 `elec_price`）
  - `elec_cost = elec_usage × elec_price`
  - `water_usage = water_current - water_last`
  - `water_cost = water_usage × water_price`
  - `total_cost = rent_cost + elec_cost + water_cost`

### 4.3 标记收费（PUT `/api/bills/:id/pay`）

- 将 `is_paid` 设为 `1`，`paid_at` 设为当前时间
- 之后账单变为**只读**（不可修改读数、不可取消编辑）

### 4.4 取消收费（PUT `/api/bills/:id/unpay`）

- 将 `is_paid` 重置为 `0`，`paid_at` 置空
- 账单恢复可编辑状态

---

## 五、API 接口及数据流

### 5.1 响应格式

所有接口统一返回：

```json
{
  "code": 0,
  "msg": "success",
  "data": { ... }
}
```

### 5.2 接口清单

| 方法 | 路径 | 说明 | 请求参数 | 响应 data |
|------|------|------|---------|-----------|
| GET | `/api/rooms` | 房间列表 | `rented`（0/1 可选）、`floor`（楼层） | `Room[]` |
| PUT | `/api/rooms/:room_no` | 修改房间 | `{ rent_price, elec_price, water_price, is_rented }` | `Room` |
| POST | `/api/bills/generate` | 生成账单 | `{ year, month }` | 生成记录数 |
| GET | `/api/bills/list` | 账单列表 | `year, month, room_no, is_paid` | `Bill[]`（含房间单价信息） |
| PUT | `/api/bills/:id` | 录入读数 | `{ elec_current, water_current }` | `Bill` |
| PUT | `/api/bills/:id/pay` | 标记收费 | — | `Bill` |
| PUT | `/api/bills/:id/unpay` | 取消收费 | — | `Bill` |
| GET | `/api/bills/:id/receipt` | 获取收据 | — | 收据数据对象 |

### 5.3 收据数据对象

```json
{
  "room_no": "201",
  "year": 2026,
  "month": 6,
  "elec_last": 1500.0,
  "elec_current": 1650.0,
  "elec_usage": 150.0,
  "elec_cost": 127.5,
  "elec_price": 0.85,
  "water_last": 88.0,
  "water_current": 96.0,
  "water_usage": 8.0,
  "water_cost": 36.0,
  "water_price": 4.5,
  "rent_price": 1200.0,
  "rent_cost": 1200.0,
  "total_cost": 1363.5
}
```

---

## 六、前端数据流

```
User Action → Vue Component → API Call → Axios → /api/proxy → Go Handler → Service → Repository → MySQL
                                                                                           │
                                                                                    GORM (自动建表)
```

### 前端路由与页面

| 路由 | 页面 | 数据来源 |
|------|------|---------|
| `/` | 房间管理 | `getRooms()`、`updateRoom()` |
| `/bills` | 月度账单 | `getBillList()`、`generateBills()`、`updateBill()`、`payBill()`、`unpayBill()` |
| `/history` | 历史账单 | `getBillList()` |
| `/dashboard` | 数据大屏 | `getRooms()` + `getBillList()` 聚合 |
| 弹窗 | 收据 | `getReceipt()` → html2canvas 截图下载 |

---

## 七、数据约束总结

| 约束 | 说明 |
|------|------|
| 房间号唯一 | `room_no` UNIQUE |
| 每房每月一账单 | `room_no + year + month` 复合唯一索引 |
| 已收费不可编辑 | 录入读数接口在校验 `is_paid` |
| 费用快照 | 生成账单时读取 `rent_price`，之后不随房间调价变化 |
| 水电费延迟计算 | 新建账单时水电费为 NULL，录入读数后自动计算 |
| 删除保护 | 无删除接口，数据仅修改状态 |
