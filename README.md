# Go линтер для проверки сообщений логирования

Линтер `loglint` реализован на базе пакета `go/analysis` и проверяет текст логов
в коде Go (1.22+), поддерживая логгеры:

- `log/slog`
- `go.uber.org/zap`

Он встроен в `go vet` как кастомный `vettool` и может подключаться к `golangci-lint`
как плагин. Ниже описаны требования, архитектура и сценарии применения.

### Требования

- Go 1.22+
- Для интеграции с `golangci-lint` - установленный `golangci-lint` совместимой версии

### Быстрый старт

1. Клонируйте репозиторий и перейдите в его корень:

   ```bash
   git clone https://github.com/wurlinney/go-log-linter.git
   cd go-log-linter
   ```

2. Соберите бинарный файл линтера:

   ```bash
   go build -o loglint ./cmd/loglint
   ```

3. Прогоните линтер как `vettool` по самому себе:

   ```bash
   go vet -vettool="$(pwd)/loglint" ./...
   ```

Линтер выведет diagnostics в формате `go vet` (`путь/к/файлу.go:строка: сообщение`).

### Применение к вашему проекту (go vet)

Предположим, что ваш проект находится в каталоге `my-service`, а собранный бинарник
`loglint` лежит в PATH или рядом.

1. Соберите `loglint` (один раз, можно из этого проекта или через `go install`):

   ```bash
   go build -o loglint ./cmd/loglint
   ```

2. В корне вашего проекта запустите:

   ```bash
   go vet -vettool="/полный/путь/до/loglint" ./...
   ```

3. Исправляйте сообщения вида:

   ```text
   internal/server/logs.go:42: log message should start with lowercase letter
   ```

Таким образом `loglint` интегрируется в стандартный рабочий процесс `go vet`.

### Правила

- **lowercase** - сообщение должно начинаться со строчной буквы.
- **english** - в сообщении допускаются только английские (ASCII) символы.
- **symbols** - запрещены спецсимволы и эмодзи, разрешены буквы, цифры и пробелы.
- **sensitive** - эвристически ищет потенциально чувствительные данные (password, token, api_key, jwt, credential, cookie, session и другие).

### Интеграция с golangci-lint

Линтер можно подключить как плагин для golangci-lint.

- **Сборка плагина**:

```bash
go build -buildmode=plugin -o loglint.so ./internal/golangci
```

- **Пример фрагмента `.golangci.yml`**:

```yaml
linters-settings:
  govet:
    enable:
      - loglint

linters:
  enable:
    - govet

run:
  timeout: 5m
```
Практический сценарий применения к вашему проекту:

1. Соберите плагин `loglint.so` (как выше).
2. Убедитесь, что golangci-lint видит `.so` (например, положите его рядом с бинарником `golangci-lint` или настройте путь согласно документации вашей версии).
3. Добавьте показанный фрагмент в `.golangci.yml` вашего проекта.
4. Запускайте:

   ```bash
   golangci-lint run ./...
   ```

Правила `loglint` будут выполняться как часть `govet` и падать вместе с общим прогоном линтеров.

### Описание функционала

- **Поддерживаемые логгеры**:
  - `log/slog` (`slog.Info`, `slog.Error`, логгеры на базе `slog.Logger`),
  - `go.uber.org/zap` (`zap.L().Info`, `zap.L().Error`, переменные-логгеры).
- **Результат работы** - diagnostics уровня `go vet`/`golangci-lint` с сообщениями правил.

### Архитектура

- internal/core - доменные типы (LogEntry, Rule, Violation).
- internal/extractors - извлечение лог вызовов из AST (slog, zap).
- internal/rules - правила проверки сообщений.
- internal/engine - пайплайн CallExpr -> LogEntry -> Rules -> Violations.
- internal/analyzer - analysis.Analyzer, обход AST и репортинг.
- internal/report - преобразование Violation в analysis.Diagnostic.
- internal/golangci - точка интеграции с golangci-lint.

Примерная структура проекта:

```text
go-log-linter/
├── cmd/
│   └── loglint/
│       └── main.go           # entrypoint для go vet (singlechecker)
├── internal/
│   ├── analyzer/             # анализатор на базе go/analysis
│   ├── core/                 # LogEntry, Rule, Violation
│   ├── engine/               # пайплайн extractors + rules
│   ├── extractors/           # парсинг slog/zap вызовов
│   ├── inspectors/           # вспомогательные проверки текста/AST
│   ├── rules/                # lowercase, english, symbols, sensitive
│   ├── report/               # преобразование в analysis.Diagnostic
│   ├── config/               # загрузка loglint.json
│   ├── golangci/             # экспорт Analyzers для golangci-lint
│   └── tools/                # фиксация зависимостей инструментов
├── testdata/
│   └── analyzer/             # сценарии для analysistest
├── .github/
│   └── workflows/ci.yml      # CI для сборки и тестов
├── go.mod
├── go.sum
└── README.md
```

### Бонусные задания

В проекте реализованы все бонусные задачи:

1. **Конфигурация правил**: включение/отключение правил и настройка `sensitive` через конфигурационный файл `loglint.json` (см. ниже).
2. **SuggestedFixes*: для правила `lowercase` линтер предлагает безопасный автофикc - первая буква строкового литерала переводится в нижний регистр.
3. **Кастомные паттерны для чувствительных данных**: в конфиге можно указать свои ключевые слова для правила `sensitive` (`sensitive.keywords`), которые будут учитываться вместе со встроенными.
4. **CI/CD**: в репозитории есть конфигурация GitHub Actions для автоматической сборки и прогона тестов при пуше (`.github/workflows/ci.yml`).

Для авто-исправлений важно, чтобы они были детерминированными и безопасными, поэтому они включены не для всех правил:

- `lowercase` - автофикс включен (безопасно заменить первую букву сообщения на строчную).
- `english` - автофикс не реализован (невозможно автоматически и однозначно перевести сообщение на английский).
- `symbols` - автофикс умышленно не добавлен, так как простое удаление символов может ломать смысл сообщения.
- `sensitive` - автофикс не реализован, потому что нет универсального и безопасного способа переписать потенциально чувствительные сообщения.

### Конфигурация

Линтер поддерживает простую JSON конфигурацию через файл loglint.json
в корне проекта (или путь из переменной окружения LOGLINT_CONFIG).

- **Пример loglint.json**:

```json
{
  "rules": {
    "lowercase": true,
    "english": true,
    "symbols": true,
    "sensitive": true
  },
  "sensitive": {
    "keywords": ["mysecret", "internal_token"]
  }
}
```

Если файл отсутствует, используются значения по умолчанию.

### Тестирование

- **Запуск всех тестов** в этом репозитории:

  ```bash
  go test ./...
  ```

  Это включает:

  - unit-тесты правил (`internal/rules/*_test.go`),
  - unit-тесты экстракторов (`internal/extractors/*_test.go`),
  - интеграционные тесты анализатора через `analysistest` и `testdata/analyzer/*`.

- **Проверка только анализатора** (актуально при изменении пайплайна/AST‑логики):

  ```bash
  go test ./internal/analyzer
  ```

При изменении структуры `testdata/analyzer/*` не забывайте обновлять ожидаемые
комментарии `// want "..."`, чтобы интеграционные тесты оставались валидными.

### Примеры

- Неверные сообщения:

```go
slog.Info("Starting server")              // нарушение lowercase
slog.Info("запуск сервера")               // нарушение english
slog.Info("server started!!!")            // нарушение symbols
slog.Info("user password: " + password)   // нарушение sensitive
```

- Корректные сообщения:

```go
slog.Info("starting server")
slog.Info("starting server on port 8080")
slog.Info("user authenticated successfully")
```
