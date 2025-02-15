# Avito Store
Этот репозиторий содержит код сервиса, необходимый для его развертывания `docker-compose.yaml`, а также инструкции по установке, запуску и тестированию.


## Установка Docker и Docker Compose

### 1. Установка Docker
#### macOS (через Homebrew):
```
brew install --cask docker
```

После установки запустите Docker Desktop и убедитесь, что он работает.

#### Ubuntu:

```
sudo apt update
sudo apt install -y docker.io
sudo systemctl enable --now docker
```


#### Windows:
Скачайте и установите [Docker Desktop](https://www.docker.com/products/docker-desktop).

---

### 2. Установка Docker Compose
На современных версиях Docker (>= 20.10) docker-compose уже встроен, но если требуется отдельная установка:

#### macOS и Linux:

```
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```


#### Windows:
Docker Compose уже встроен в Docker Desktop.

---

## Запуск сервиса

1. Клонируйте репозиторий:

```
git clone https://github.com/reden1k/avito-tech-winter-2025.git
cd avito-tech-winter-2025
```


2. Запустите сервис с пересборкой образов:

docker-compose up --build

После успешного запуска сервис будет доступен по 8080 порту.

---

## Тестирование

### 1. e2e тесты
Запуск тестов:

```
go test ./test/e2e
```


### 2. Нагрузочное тестирование через k6
Установите k6, если он не установлен:

```
brew install k6  # для macOS
sudo apt install k6  # для Ubuntu
choco install k6  # для Windows (Chocolatey)
```


Запустите тест:

```
cd avito-tech-winter-2025
docker-compose up --build
k6 run load_test.js
```


Тест проверит:
- RPS (запросов в секунду)
- SLI по времени ответа
- SLI по статусу ответа

После выполнения будут выведены метрики, которые можно использовать для анализа производительности.


### Результаты нагрузочного тестирования
![изображение](https://github.com/user-attachments/assets/699f5f9a-b318-4801-9bc0-2ba94ee979c3)
