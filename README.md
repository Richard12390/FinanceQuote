# Financial Quote Streaming Platform

本專案是一套以 **Event-Driven Architecture** 為核心的即時行情系統，  
整合 **Go、NATS、Java Spring Boot 與 Vue**，用於蒐集、處理並即時推送  
**股票、ETF、虛擬貨幣** 的秒級與日級行情資料。

系統重點在於 **高併發行情擷取、跨市場資料正規化、服務解耦與即時推播**，  
並以 **DDD（Domain-Driven Design）** 與 **Hexagonal Architecture** 組織程式碼。

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

此設計使下游服務可以依需求訂閱特定市場或資產類型，  
避免耦合於特定資料來源。

---

### Java – Backend Service

Java 後端使用 **Spring Boot**，負責行情處理與對外服務：

- 訂閱 NATS 行情事件並反序列化為 `QuoteNorm`
- 透過 Domain Service 補足交易所、資產分類等領域語意
- 日級行情資料使用 **MyBatis** 寫入 **MySQL**
- 秒級行情資料透過 **WebSocket** 即時推送前端
- 提供 REST API，支援：
  - 行情查詢
  - 使用者操作
  - JWT 登入驗證

---

### Frontend – Vue

前端使用 **Vue 3**，提供即時資料呈現與使用者操作介面：

- 使用者登入（JWT 驗證）
- 即時行情總覽
- 使用者可將標的加入或移除「我的最愛（Favorites）」

---

## Architecture Style

本專案在設計上採用：

- Event-Driven Architecture
- Domain-Driven Design（DDD）
- Hexagonal Architecture

系統清楚劃分 **Domain / Application / Adapter / Infrastructure**，  
使核心業務邏輯不直接依賴外部系統，提升可維護性與可擴充性。

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

---

##Tech Stack

-Go：高併發行情擷取與事件發佈
-NATS：事件訊息中樞
- Java 17 / Spring Boot：後端服務
- MyBatis / MySQL：資料存取
- JWT：登入驗證
- WebSocket：即時資料推送
- Vue 3：前端介面
- Docker Compose：部署與服務管理
- Zap + Lumberjack：結構化日誌與輪轉
