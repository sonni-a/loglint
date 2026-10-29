# LogLint

[![CI](https://github.com/sonni-a/loglint/actions/workflows/ci.yml/badge.svg)](https://github.com/sonni-a/loglint/actions/workflows/ci.yml)

Линтер на Go, совместимый с golangci-lint, анализирующий лог-записи в коде и проверяющий их соответствие установленным правилам.
Поддерживаются `log/slog`, стандартный `log` и логгеры с типом, содержащим `logger` или `zap` (в том числе zap-совместимые обёртки).

## Функционал

* **CheckLowercase**: лог-сообщения должны начинаться со строчной буквы.
* **CheckSpecialChars**: лог-сообщения не должны содержать спецсимволы или эмодзи.
* **CheckEnglish**: лог-сообщения должны быть только на английском языке (ASCII).
* **CheckSensitive**: лог-сообщения не должны содержать потенциально чувствительные данные.

Для одного сообщения линтер сообщает обо всех нарушенных правилах. Autofix применяет исправления по цепочке: sensitive → english → special → style.

## Инструкция по запуску

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
5. Собрать standalone-версию линтера:
   ```bash
   go build -o loglint ./cmd/loglint
   ```
6. Чтобы собрать custom-gcl, нужен golangci-lint v2. В корне проекта:
   ```bash
   golangci-lint custom
   ```
7. Запустить анализ на другом проекте, скопировав `custom-gcl` (на Windows — `custom-gcl.exe`) и `.golangci.yml` в корень проекта:
   ```bash
   ./custom-gcl run
   ```
8. Запустить линтер с autofix:
   ```bash
   ./custom-gcl run --fix
   ```

Настройки линтера можно изменить в `.golangci.yml` — после этого пересобирать бинарник не нужно.

Примеры использования:

![Img 1](https://i127.fastpic.org/big/2026/0314/c2/7c673096db7d82aa7a368447ec2102c2.jpeg)

![Img2](https://i127.fastpic.org/big/2026/0314/aa/9d82ad05a001c5c63b42b38319042eaa.png)

После `--fix` останутся только те нарушения, для которых нет автоматического исправления (например, сообщения не на английском языке).
