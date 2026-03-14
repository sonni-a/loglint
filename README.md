# LogLint

Линтер на Go, совместимый с golangci-lint, анализирующий лог-записи в коде и проверяющий их соответствие  установленным правилам.
Линтер поддерживает slog и zap.

## Функционал

Анализатор поддерживает следующие проверки, которыми можно управлять через файл конфигурации `.custom-gcl.yml`:

* **CheckLowercase**: лог-сообщения должны начинаться со строчной буквы.
* **CheckSpecialChars**: лог-сообщения не должны содержать спецсимволы или эмодзи.
* **CheckEnglish**: лог-сообщения должны быть только на английском языке.
* **CheckSensitive**: лог-сообщения не должны содержать потенциально чувствительные данные.

В sensitive-words  можно через запятую перечислить паттерны для проверки чувствительных данных.
В флажках enable-{правило} можно поставить false, чтобы линтер их игнорировал.

## Инструкция по конфигурации и запуску

1. Склонировать репозиторий
   ```bash
   git clone https://github.com/sonni-a/loglint.git
   ```
2. Перейти в папку с проектом loglint
   ```bash
   cd loglint
   ``` 
3. Скачать зависимости
   ```bash
   go mod tidy
   ```
4. Выполнить тестирование:
  ```bash
  go test -v ./...
  ```
5. При желании изменить в файле golangci.yml изменить найстройки правил
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
          sensitive-words: "password,api_key,token,secret,credit_card" # тут можно изменить sensitive-words
          enable-sensitive: true
          enable-style: true
          enable-english: true
          enable-special: true
   ```
6. Для получения .exe файла линтера выполните:
  ```bash
  go build -o loglint.exe ./cmd/loglint/main.go
  ```
7. Чтобы запустить линтер на другом проекте, перейдите в корневую папку этого проекта и выполните:
  ```bash
  [путь к loglint]/loglint/loglint.exe ./...
  ``` 
8. Чтобы запустить линтер с Suggested Autofixes:
  ```bash
  [путь к loglint]/loglint/loglint.exe -fix ./...
  ``` 
