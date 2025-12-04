# Financial Quote Streaming Platform

本專案是一套即時行情系統，  
整合 **Go、NATS、Java Spring Boot、RabbitMQ 與 Vue**，  
用於蒐集、處理並推送 **股票、ETF、虛擬貨幣** 的秒級與日級行情資料。

系統以 **事件驅動的通訊模式** 串接不同服務：
- Go 負責行情擷取與事件發佈
- Java 負責事件處理、資料落庫與即時推播
- 前端透過 WebSocket 即時接收行情更新

在程式碼組織上，本專案**參考 DDD（Domain-Driven Design）與 Hexagonal Architecture 的分層概念**，嘗試將**核心業務邏輯**、**資料模型**及**基礎設施實作**分離，達到職責清楚、機制可替換的設計目標。

---

## System Overview

### Go – Market Data Producer

Go 服務負責行情資料的高併發擷取與事件發佈：

- 依 YAML pipeline 設定，定時抓取秒級 / 日級行情
- 擷取資料來源包含：
  - Yahoo Finance
  - TWSE
  - Binance
- 不同來源資料會被轉換並統一正規化為 `QuoteNorm` 結構
- 行情事件依「資料來源」與「資產類型」發佈至 NATS，例如：
  - `quotes.yahoo.stock`
  - `quotes.twse.etf`
  - `quotes.binance.crypto`

此設計使下游服務可以依需求訂閱特定市場或資產類型，避免耦合於特定資料來源。

---

### Java – Backend Service

Java 後端使用 **Spring Boot**，負責行情處理與對外服務：

- 訂閱 NATS 行情事件並反序列化為 `QuoteNorm`
- 透過 Domain Service 補足交易所、資產分類等領域語意
- 日級行情資料使用 **MyBatis** 寫入 **MySQL**
- 提供 REST API，支援行情查詢與 **JWT** 登入驗證
- 秒級行情資料由 Java 服務轉交至 **RabbitMQ**，透過 Spring 的 **STOMP Broker Relay** 機制，將訊息經由 **WebSocket** 即時推送至前端瀏覽器
- 秒級行情資料在後端使用 **RabbitMQ Stream** 作為高頻資料通道，提供資料緩衝與可重播能力，以因應連續且密集的行情更新

---

### Frontend – Vue

前端使用 **Vue 3**，提供即時資料呈現與使用者操作介面：

- 使用者登入（JWT 驗證）
- 即時行情總覽
- 使用者可將標的加入或移除「我的最愛（Favorites）」

---

## Data Flow

### Go Data Pipeline

1. 啟動時讀取 YAML pipeline 設定
2. 初始化 NATS、Logger（Zap + Lumberjack）與各資產 Fetcher
3. 使用 `context` 統一管理取消與逾時
4. 依設定進行秒級 / 日級排程輪詢
5. 使用雙層併發控制：
   - Worker Pool 控制任務數量
   - Semaphore 限制對外 API 請求量
6. 將原始資料轉換為 `QuoteNorm`
7. 發佈事件至 NATS

---

### Java Data Processing

1. 訂閱 NATS subjects 接收行情事件
2. 解析並驗證 `QuoteNorm`
3. 補足交易所與資產分類等領域資訊
4. 日級行情落庫至 MySQL
5. 秒級行情即時推送至前端

---

## Deployment

本專案使用 **Docker Compose** 管理所有服務。

```bash
cd deploy
docker compose up -d

服務啟動後可於瀏覽器存取前端介面：http://localhost:8030
```

---

## Tech Stack

- **Go**：行情資料擷取與事件發佈，使用 worker pool 與 semaphore 控制多標的併發請求
- **NATS**：後端服務間的事件通訊中樞（Go → Java）
- **Java 17 / Spring Boot**：後端服務，負責事件處理、資料落庫與即時推播
- **MyBatis / MySQL**：日級行情與基礎資料的持久化存取
- **RabbitMQ**：即時推播用的訊息代理，負責 routing 與多使用者 fan-out
- **RabbitMQ Stream**：高頻行情資料通道，提供緩衝與可重播能力
- **STOMP（over WebSocket）**：即時推播協議，定義前端訂閱與訊息傳遞語意
- **WebSocket**：與瀏覽器保持即時雙向連線
- **JWT**：登入驗證與 API 存取控制
- **Vue 3**：前端即時行情顯示與使用者操作介面
- **Docker Compose**：多服務開發與部署管理
- **Zap + Lumberjack**：結構化日誌與每日輪轉

---
