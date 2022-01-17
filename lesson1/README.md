
# Запуск
```
docker-compose -f docker-compose.yml up -d
```
Чтобы настроить Grafana, зайдите в личный кабинет(логин - admin, пароль - Password), Settings -> Data Sources -> Add data source -> Prometheus и в поле URL вставьте ```http://prometheus:9090``` .