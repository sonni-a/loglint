# LogLint

Линтер на Go, совместимый с golangci-lint, анализирующий лог-записи в коде и проверяющий их соответствие  установленным правилам.
Линтер поддерживает slog и zap.

## Функционал и конфигурация

Анализатор поддерживает следующие проверки, которыми можно управлять через файл конфигурации `.custom-gcl.yml`:

* **CheckLowercase**: лог-сообщения должны начинаться со строчной буквы.
* **CheckSpecialChars**: лог-сообщения не должны содержать спецсимволы или эмодзи.
* **CheckEnglish**: лог-сообщения должны быть только на английском языке.
* **CheckSensitive**: лог-сообщения не должны содержать потенциально чувствительные данные.

### Настройка правил и кастомных паттернов через файл .golangci.yml
```yaml
linters:
  disable-all: true
  enable:
    - loglint

linters-settings:
  custom:
    loglint:
      type: "module"
      description: "Linter for checking log messages style and security"
      settings:
        sensitive-words: "password,api_key,token,secret,credit_card"
        enable-sensitive: true 
        enable-style: true
        enable-english: true
        enable-special: true
```
В sensitive-words  можно через запятую перечислить паттерны для проверки чувствительных данных.
В флажках enable-{правило} можно поставить false, чтобы линтер их игнорировал.

## Инструкция по запуску
Для получения .exe файла линтера выполните:
```bash
go build -o loglint ./cmd/loglint/main.go
```
Запуск линтера:
```bash
./loglint ./...
```
Запуск линтера с Suggested Autofixes:
```bash
./loglint -fix ./...
```

## Тестирование
Для запуска unit- и интеграционных тестов:
```bash
go test -v ./...
```
