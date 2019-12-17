
### Configuration

```bash
cat <<EOF > config.toml
[sender]
max_workers = 1  # максимальное кол-во обработчиков
unmarshaler = "json"  # декодер сообщений в кафке ( `json` / `proto` )

[kafka]
brokers = "127.0.0.1:9092"  # список брокеров через запятую
topic = "messages"  # топик с задачами на отпраку сообщений
client_id = "srv-email"  # идентификатор клиента

[smtp]
server_addr = "smtp.example.com:465"  # адрес smtp сервера (хост и порт)
username = "sender@example.com"  # имя пользователя
password = "smtp-password"  # пароль
EOF
```

### Running

```bash
go run main.go
```

### Send messages to queue

```bash
brew install kafkacat fortune

echo '{"from":"sender@example.com","recipients":["kazhuravlev@fastmail.com"],"headers":{},"subject":"New Fortune!","body":"'`fortune | base64`'"}' | kafkacat -b 127.0.0.1:9092 -P -t messages
```
