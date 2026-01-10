# Ваш первый проект с нуля

> **Быстрый старт** - создайте свой собственный MT5 торговый проект за 10 минут, используя только Go модули

---

## Для кого этот гайд?

Этот документ предназначен для тех, кто хочет:

- **Быстро начать** писать код для MT5 в своем проекте на Go
- **Не клонировать** весь репозиторий GoMT5
- **Создать проект с нуля** и подключить минимальные зависимости
- **Написать первый метод** и увидеть результат немедленно

**Разница между этим гайдом и GETTING_STARTED.md:**

| Getting Started | Your First Project (этот гайд) |
|----------------|--------------------------------|
| Клонируете готовый репозиторий | Создаете проект с нуля |
| Изучаете архитектуру и примеры | Сразу пишете работающий код |
| Долгий путь обучения | Быстрый результат |
| Для глубокого погружения | Для быстрого старта |

> После того как вы пройдете этот гайд и получите первый результат, рекомендуем изучить [GETTING_STARTED.md](./GETTING_STARTED.md) для понимания полной архитектуры SDK.

---

## Что мы будем делать?

В этом гайде мы создадим минималистичный проект, который:

1. Подключится к MT5 терминалу через gRPC шлюз
2. Получит баланс счета
3. Выведет результат в консоль

**Это займет 10 минут и требует минимум кода.**

---

## Шаг 1: Установите Go 1.23 или выше

Если у вас еще не установлен Go:

**Скачайте и установите:**

- [Go Download](https://go.dev/dl/)

**Проверьте установку:**

```bash
go version
# Должно показать: go version go1.23.x или выше
```

---

## Шаг 2: Создайте новый Go проект

Откройте терминал (командную строку) и выполните:

```bash
# Создаем папку для проекта
mkdir MyMT5Project
cd MyMT5Project

# Инициализируем Go модуль
go mod init mymt5project
```

**Что произошло:**

- Создана папка `MyMT5Project`
- Внутри создан файл `go.mod` - манифест вашего проекта
- Теперь можно добавлять зависимости

---

## Шаг 3: Установите необходимые пакеты

Это самый важный шаг - устанавливаем пакеты для работы с MT5:

```bash
# Основной MT5 клиент
go get github.com/MetaRPC/GoMT5/mt5

# Protobuf определения MT5 API
go get git.mtapi.io/root/mrpc-proto/mt5/libraries/go

# UUID для идентификаторов
go get github.com/google/uuid

# gRPC и protobuf
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

**Что включают эти пакеты:**

- `github.com/MetaRPC/GoMT5/mt5` - MT5Account для работы с MT5
- `git.mtapi.io/root/mrpc-proto` - Protocol Buffers схемы для gRPC API
- `github.com/google/uuid` - генерация UUID для сессий
- `google.golang.org/grpc` - gRPC клиент
- `google.golang.org/protobuf` - protobuf runtime

> **Важно:** Эти пакеты - это ВСЁ что вам нужно для работы с MT5 через Go. Никаких дополнительных файлов клонировать не требуется.

---

## Шаг 4: Создайте файл конфигурации config.json

Создайте файл `config.json` в корне проекта:

```json
{
  "User": 591129415,
  "Password": "IpoHj17tYu67@",
  "MtCluster": "FxPro-MT5 Demo",
  "GrpcServer": "mt5.mrpc.pro:443",
  "TestSymbol": "EURUSD",
  "ConnectTimeout": 120
}
```

**Объяснение параметров:**

| Параметр | Описание | Пример |
|----------|----------|--------|
| **User** | Номер вашего MT5 счета (логин) | `591129415` |
| **Password** | Мастер-пароль от MT5 счета | `"IpoHj17tYu67@"` |
| **MtCluster** | Название кластера вашего брокера | `"FxPro-MT5 Demo"` |
| **GrpcServer** | Адрес gRPC шлюза (host:port) | `"mt5.mrpc.pro:443"` |
| **TestSymbol** | Торговый символ по умолчанию | `"EURUSD"` |
| **ConnectTimeout** | Таймаут подключения в секундах | `120` |

**Замените:**

- `User`, `Password`, `MtCluster` - на данные вашего MT5 демо-счета
- `GrpcServer` - оставьте как есть (это адрес публичного шлюза MetaRPC)

> **Нет MT5 аккаунта?** Создайте демо-счет у любого брокера с MT5 терминалом (FxPro, Alpari, RoboForex и т.д.)

---

## Шаг 5: Создайте структуру для конфигурации

Создайте файл `config.go` в корне проекта:

```go
package main

import (
	"encoding/json"
	"os"
)

// Config содержит настройки подключения к MT5
type Config struct {
	User           uint64 `json:"User"`
	Password       string `json:"Password"`
	MtCluster      string `json:"MtCluster"`
	GrpcServer     string `json:"GrpcServer"`
	TestSymbol     string `json:"TestSymbol"`
	ConnectTimeout int    `json:"ConnectTimeout"`
}

// LoadConfig загружает конфигурацию из JSON файла
func LoadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
```

---

## Шаг 6: Напишите код для подключения и получения баланса

Создайте файл `main.go` в корне проекта:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "git.mtapi.io/root/mrpc-proto/mt5/libraries/go"
	"github.com/MetaRPC/GoMT5/mt5"
	"github.com/google/uuid"
)

func main() {
	// ============================================================================
	// КОНФИГУРАЦИЯ - Загружаем настройки из config.json
	// ============================================================================

	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Println("=== MT5 Connection Configuration ===")
	fmt.Printf("User: %d\n", config.User)
	fmt.Printf("Cluster: %s\n", config.MtCluster)
	fmt.Printf("gRPC Server: %s\n", config.GrpcServer)
	fmt.Printf("Symbol: %s\n", config.TestSymbol)
	fmt.Println("====================================\n")

	// ============================================================================
	// ПОДКЛЮЧЕНИЕ - Создаем MT5Account и подключаемся к MT5
	// ============================================================================

	fmt.Println("Connecting to MT5 gateway...")

	// Создаем контекст
	ctx := context.Background()

	// Создаем MT5Account - это главный объект для работы с MT5
	account, err := mt5.NewMT5Account(
		config.User,
		config.Password,
		config.GrpcServer,
		uuid.New(), // Генерируем уникальный ID сессии
	)
	if err != nil {
		log.Fatalf("Failed to create MT5Account: %v", err)
	}
	defer account.Close()

	// Подключаемся к MT5 терминалу через ConnectEx
	baseSymbol := config.TestSymbol
	timeoutSec := int32(config.ConnectTimeout)
	connectReq := &pb.ConnectExRequest{
		User:                                   config.User,
		Password:                               config.Password,
		MtClusterName:                          config.MtCluster,
		BaseChartSymbol:                        &baseSymbol,
		TerminalReadinessWaitingTimeoutSeconds: &timeoutSec,
	}

	connectData, err := account.ConnectEx(ctx, connectReq)
	if err != nil {
		log.Fatalf("ConnectEx failed: %v", err)
	}

	// Обновляем ID сессии из ответа
	account.Id = uuid.MustParse(connectData.TerminalInstanceGuid)

	fmt.Println("✓ Connected successfully!")
	fmt.Printf("Instance ID: %s\n\n", connectData.TerminalInstanceGuid)

	// ============================================================================
	// ПОЛУЧЕНИЕ БАЛАНСА - Вызываем метод AccountSummary
	// ============================================================================

	fmt.Println("Fetching account balance...")

	accountSummaryReq := &pb.AccountSummaryRequest{}

	accountData, err := account.AccountSummary(ctx, accountSummaryReq)
	if err != nil {
		log.Fatalf("AccountSummary failed: %v", err)
	}

	// Выводим информацию о счете
	fmt.Println("=== Account Information ===")
	fmt.Printf("Login: %d\n", accountData.AccountLogin)
	fmt.Printf("Balance: %.2f\n", accountData.AccountBalance)
	fmt.Printf("Equity: %.2f\n", accountData.AccountEquity)
	fmt.Printf("Margin: %.2f\n", accountData.AccountMargin)
	fmt.Printf("Free Margin: %.2f\n", accountData.AccountFreeMargin)
	fmt.Printf("Currency: %s\n", accountData.AccountCurrency)
	fmt.Printf("Leverage: 1:%d\n", accountData.AccountLeverage)
	fmt.Printf("Company: %s\n", accountData.AccountCompanyName)
	fmt.Printf("User Name: %s\n", accountData.AccountUserName)

	if accountData.ServerTime != nil {
		serverTime := accountData.ServerTime.AsTime()
		fmt.Printf("Server Time: %s\n", serverTime.Format(time.RFC3339))
	}

	fmt.Println("===========================\n")

	fmt.Println("✓ Success! Your first MT5 connection is complete.")

	// ============================================================================
	// ОТКЛЮЧЕНИЕ - Закрываем соединение
	// ============================================================================

	disconnectReq := &pb.DisconnectRequest{}
	_, err = account.Disconnect(ctx, disconnectReq)
	if err != nil {
		log.Printf("Disconnect warning: %v", err)
	} else {
		fmt.Println("✓ Disconnected from MT5.")
	}
}
```

---

## Шаг 7: Запустите проект

Сохраните все файлы и выполните:

```bash
go run .
```

**Или соберите и запустите:**

```bash
go build -o mymt5project.exe
./mymt5project.exe
```

**Ожидаемый результат:**

```
=== MT5 Connection Configuration ===
User: 591129415
Cluster: FxPro-MT5 Demo
gRPC Server: mt5.mrpc.pro:443
Symbol: EURUSD
====================================

Connecting to MT5 gateway...
✓ Connected successfully!
Instance ID: 12345678-90ab-cdef-1234-567890abcdef

Fetching account balance...
=== Account Information ===
Login: 591129415
Balance: 10000.00
Equity: 10000.00
Margin: 0.00
Free Margin: 10000.00
Currency: USD
Leverage: 1:100
Company: FxPro Financial Services Ltd
User Name: Demo User
Server Time: 2025-01-06T15:30:45Z
===========================

✓ Success! Your first MT5 connection is complete.
✓ Disconnected from MT5.
```

---

## Поздравляем! Вы сделали это!

Вы только что:

✅ Создали новый Go проект с нуля
✅ Подключили Go модули для работы с MT5
✅ Настроили конфигурацию подключения
✅ Подключились к MT5 терминалу через gRPC
✅ Получили баланс счета программно

**Это был низкоуровневый (Low-Level) подход** с прямым использованием `MT5Account` и gRPC.

---

## Что дальше?

Теперь, когда у вас есть рабочий проект, вы можете:

### 1. Изучить полную архитектуру SDK

Прочитайте [GETTING_STARTED.md](./GETTING_STARTED.md) чтобы узнать о:

- **MT5Account** (Low-Level) - то что вы только что использовали
- **MT5Service** (Wrappers) - удобные обертки над MT5Account
- **MT5Sugar** (High-Level) - синтаксический сахар для быстрой разработки

### 2. Добавить больше функциональности

**Примеры того что можно сделать:**

```go
// Получить все открытые позиции
openedOrdersReq := &pb.OpenedOrdersRequest{
    InputSortMode: pb.BMT5_ENUM_OPENED_ORDER_SORT_TYPE_SORT_BY_TIME_OPEN,
}
openedData, err := account.OpenedOrders(ctx, openedOrdersReq)

// Открыть рыночный ордер на покупку (через MT5Sugar)
sugar := mt5sugar.NewMT5Sugar(service, account.Id.String())
ticket, err := sugar.BuyMarket("EURUSD", 0.01)
if err == nil {
    fmt.Printf("Order opened: #%d\n", ticket)
}

// Получить котировки в реальном времени (streaming)
tickReq := &pb.OnSymbolTickRequest{SymbolNames: []string{"EURUSD"}}
dataChan, errChan := account.OnSymbolTick(ctx, tickReq)

for {
    select {
    case data := <-dataChan:
        tick := data.SymbolTick
        fmt.Printf("Bid: %.5f, Ask: %.5f\n", tick.Bid, tick.Ask)
    case err := <-errChan:
        fmt.Printf("Error: %v\n", err)
        return
    case <-ctx.Done():
        return
    }
}
```

### 3. Использовать готовые классы из репозитория

Если вы хотите использовать **MT5Service** или **MT5Sugar** в своем проекте:

1. Склонируйте репозиторий GoMT5
2. Скопируйте файлы из `examples/mt5/` в свой проект:
   - `MT5Account.go` (если хотите свою версию)
   - `MT5Service.go` (удобные обертки)
   - `MT5Sugar.go` (высокоуровневый API)
3. Используйте удобные методы высокого уровня

**Пример с MT5Sugar:**

```go
// Создаем Sugar API поверх Service
service := mt5.NewMT5Service(account)
sugar := mt5sugar.NewMT5Sugar(service, account.Id.String())

// Открыть Buy позицию
ticket, err := sugar.BuyMarket("EURUSD", 0.01)

// Закрыть все позиции по символу
err = sugar.CloseAllPositions("EURUSD")

// Получить баланс одной строкой
balance, err := sugar.GetBalance()
```

### 4. Изучить готовые примеры

В репозитории GoMT5 есть множество примеров в папке `examples/demos/`:

- `lowlevel/` - примеры работы с MT5Account
- `service/` - примеры с MT5Service
- `sugar/` - примеры с MT5Sugar (высокоуровневый API)

### 5. Прочитать дополнительные гайды

- [MT5Account API](../API_Reference/MT5Account.md) - полный справочник MT5Account
- [MT5Service API](../API_Reference/MT5Service.md) - полный справочник MT5Service
- [MT5Sugar API](../API_Reference/MT5Sugar.md) - полный справочник MT5Sugar
- [GRPC_STREAM_MANAGEMENT.md](./GRPC_STREAM_MANAGEMENT.md) - работа с потоковыми данными
- [RETURN_CODES_REFERENCE.md](./RETURN_CODES_REFERENCE.md) - коды возврата операций
- [PROTOBUF_INSPECTOR_GUIDE.md](./PROTOBUF_INSPECTOR_GUIDE.md) - инспектор protobuf структур

---

## Частые вопросы (FAQ)

### Где взять доступ к gRPC шлюзу?

В примере используется публичный шлюз MetaRPC:

```
Host: mt5.mrpc.pro
Port: 443
```

Этот шлюз доступен всем для тестирования.

> Если у вас есть вопросы по работе шлюза, посетите [GitHub Issues](https://github.com/MetaRPC/GoMT5/issues).

### Могу ли я использовать свой собственный шлюз?

Да! Если у вас есть собственная инстанция шлюза, просто измените параметр `GrpcServer` в `config.json`.

### Как получить MT5 демо-счет?

1. Скачайте MT5 терминал с сайта любого брокера (FxPro, Alpari, RoboForex)
2. Установите терминал
3. Откройте терминал и выберите "Файл" → "Открыть счет"
4. Выберите "Демо-счет" и следуйте инструкциям
5. Запишите полученные данные: логин, пароль, название сервера

### Что если я получаю ошибку подключения?

Проверьте:

1. Правильность логина/пароля/кластера в `config.json`
2. Что MT5 терминал не запущен локально (шлюз сам подключается к MT5)
3. Интернет-соединение
4. Таймаут подключения (увеличьте `ConnectTimeout` если медленный интернет)

### Нужно ли устанавливать MT5 терминал на мою машину?

**Нет!** Шлюз MetaRPC сам подключается к серверам MT5. Вам нужны только:

- Логин/пароль от MT5 счета
- Название кластера брокера
- Доступ к gRPC шлюзу

### Как работать с переменными окружения вместо config.json?

Вы можете использовать переменные окружения:

```go
import "os"

config := &Config{
    User:       uint64(os.Getenv("MT5_USER")),
    Password:   os.Getenv("MT5_PASSWORD"),
    MtCluster:  os.Getenv("MT5_CLUSTER"),
    GrpcServer: os.Getenv("MT5_GRPC_SERVER"),
    TestSymbol: os.Getenv("MT5_TEST_SYMBOL"),
}
```

Установите переменные:
```bash
export MT5_USER=591129415
export MT5_PASSWORD="IpoHj17tYu67@"
export MT5_CLUSTER="FxPro-MT5 Demo"
export MT5_GRPC_SERVER="mt5.mrpc.pro:443"
export MT5_TEST_SYMBOL="EURUSD"
```

---

## Структура вашего проекта

После завершения всех шагов ваша структура проекта должна выглядеть так:

```
MyMT5Project/
├── config.json          # Конфигурация подключения
├── config.go            # Загрузка конфигурации
├── main.go              # Главный код приложения
├── go.mod               # Go модуль с зависимостями
├── go.sum               # Контрольные суммы зависимостей
└── mymt5project.exe     # Скомпилированное приложение (после go build)
```

---

## Резюме: Что мы сделали

В этом гайде вы создали минималистичный проект, который:

1. **Использует только Go модули** - не требует клонирования репозитория
2. **Подключается к MT5** через gRPC шлюз
3. **Читает конфигурацию** из `config.json`
4. **Выполняет низкоуровневые gRPC вызовы** напрямую через `MT5Account`
5. **Получает баланс счета** и выводит в консоль

**Это основа** для любого вашего MT5 проекта на Go.

---

## Следующие шаги

Теперь вы готовы к:

- 📖 [GETTING_STARTED.md](../GETTING_STARTED.md) - Полное изучение архитектуры SDK
- 📖 [QUICK_REFERENCE.md](../QUICK_REFERENCE.md) - Быстрый справочник всех методов
- 📖 [API_REFERENCE.md](../API_REFERENCE.md) - Подробный API справочник
- 📖 [GRPC_STREAM_MANAGEMENT.md](../GRPC_STREAM_MANAGEMENT.md) - Управление потоками
- 🎯 Изучению готовых примеров в `examples/demos/`

---

## Полный пример `go.mod`

После всех установок ваш `go.mod` должен выглядеть примерно так:

```go
module mymt5project

go 1.23

require (
	git.mtapi.io/root/mrpc-proto v0.0.0-20250812093834-58b4119a2c55
	github.com/MetaRPC/GoMT5/mt5 v0.0.0
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.7
)

require (
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
)
```

---

