# LogLint

[![CI](https://github.com/sonni-a/loglint/actions/workflows/ci.yml/badge.svg)](https://github.com/sonni-a/loglint/actions/workflows/ci.yml)

Линтер для golangci-lint: проверяет текст сообщений в вызовах логгера.
Поддерживаются `log`, `log/slog`, `go.uber.org/zap` и типы с `logger` в имени.

## Функционал

* **CheckLowercase**: лог-сообщения должны начинаться со строчной буквы.
* **CheckSpecialChars**: в сообщении не должно быть `!`, `?` и эмодзи.
* **CheckEnglish**: лог-сообщения должны быть только на английском языке (ASCII).
* **CheckSensitive**: в сообщении не должно быть заданных слов вроде `password` и `token`.

Для одного сообщения линтер сообщает обо всех нарушенных правилах. Autofix применяет исправления по цепочке: sensitive → english → special → style.

## Инструкция по запуску

Нужен Go 1.27.

1. Склонировать репозиторий и перейти в него:
   ```bash
   git clone https://github.com/sonni-a/loglint.git
   cd loglint
   ```
2. Прогнать тесты:
   ```bash
   go test ./...
   ```
3. Собрать линтер и запустить его в этом модуле:
   ```bash
   go build -o loglint ./cmd/loglint
   ./loglint ./...
   ./loglint -fix ./...
   ```
   На Windows: `.\loglint ./...`

Подключить к golangci-lint v2 можно так. Команда `golangci-lint` должна быть в `PATH`.

```bash
golangci-lint custom
```

Скопируйте `custom-gcl` (на Windows — `custom-gcl.exe`) и `.golangci.yml` в корень другого Go-модуля:

```bash
./custom-gcl run
./custom-gcl run --fix
```

Правила меняются в `.golangci.yml`, пересобирать бинарник для этого не нужно.

Примеры использования:

![Img 1](https://i127.fastpic.org/big/2026/0314/c2/7c673096db7d82aa7a368447ec2102c2.jpeg)

![Img2](https://i127.fastpic.org/big/2026/0314/aa/9d82ad05a001c5c63b42b38319042eaa.png)

После `--fix` останутся только те нарушения, для которых нет автоматического исправления (например, сообщения не на английском языке).