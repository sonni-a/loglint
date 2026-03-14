# LogLint

Линтер на Go, совместимый с golangci-lint, анализирующий лог-записи в коде и проверяющий их соответствие  установленным правилам.
Линтер поддерживает slog и zap.

## Функционал

* **CheckLowercase**: лог-сообщения должны начинаться со строчной буквы.
* **CheckSpecialChars**: лог-сообщения не должны содержать спецсимволы или эмодзи.
* **CheckEnglish**: лог-сообщения должны быть только на английском языке.
* **CheckSensitive**: лог-сообщения не должны содержать потенциально чувствительные данные.

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
5. Чтобы собрать custom-gcl.exe, выполните в корне проекта:
  ```bash
  golangci-lint custom
  ```
6. Можно запустить анализ на другом проекте, скопировав перед этим custom-gcl.exe и .golangci.yml в корневую папку проекта:
  ```bash
  ./custom-gcl.exe run
  ``` 
7. Чтобы запустить линтер с Suggested Autofixes:
  ```bash
  ./custom-gcl.exe run --fix
  ```
При желании можно изменить настройки линтера в .golangci.yml и сохранить файл, заново собирать бинарник при этом не нужно!

Примеры использования:
![img](https://i127.fastpic.org/big/2026/0314/c2/7c673096db7d82aa7a368447ec2102c2.jpeg)
![img](https://i127.fastpic.org/big/2026/0314/aa/9d82ad05a001c5c63b42b38319042eaa.jpeg)
После -fix останутся только log message must be in English.
