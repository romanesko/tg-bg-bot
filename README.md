## Конфигурация бота

config.yaml:
```yaml
bot:
  token: 40257161:AAG0-5Sаfdsf9dsfSFDfdss-asd
api:
  config_url: https://somehost.com/api/bot/config.php
  token: some_test_token
  queue: 
    url: https://somehost/messages_queue
    interval: 5 

```

При запуске бот запросит `config_url`, и ожидает получить следующий конфиг:


```json
{
  "status": "complete",
  "response": {
    "start": "/menu.php",
    "buttons_header": "Выберите пункт меню ↓"   
  }
}
```    

Все url'ы могут быть как абсолютными, так и относительными.
Если не указан абсолютный, то будут подставлены протокол и хост из `config_url`

Например, при передаче относительного урла `/menu.php`, бот сделает запрос к `https://somehost.com/menu.php` ( при указаннном пути конфига `https://somehost.com/api/bot/config.php`)

## Ответы сервера

Все ответы содержат следующую структуру:

```json
{
  "status": "complete",
  "response": {
    "messages": [
      {
        "text": "Первое сообщение",
        "images": [
          "https://url-to-image",
          "https://url-to-image-two"
        ]
      },
      {
        "text": "Второе сообщение"
      },
      {
        "text": "Третье сообщение",
        "images": [
          "https://url-to-another-image"
        ]
      }
    ],
    "buttons_header": "Куда дальше?",
    "buttons": [
      {
        "label": "Пойти туда",
        "url": "/relative-url-to-next-action",
        "context" : {}
      },
      {
        "label": "Пойти сюда",
        "url": "https://absplute-url-to-next-action",
        "context" : {}
      }
    ],
    "callback" : {
      "url": "/calback-url",
      "context" : {}
    }
  }
}
```

Ответ может не содержать сообщений, а только, например, кнопки (разумно для главного меню).
Также можно переопределить `buttons_header` для каждого конкретного сообщения:

```json
{
  "status": "complete",
  "response": {
    "buttons_header": "Это главное меню! (это сообщение пропадёт вместе с кнопками)",
    "buttons": [
      {
        "label": "Пойти туда",
        "url": "/relative-url-to-next-action",
      },
      {
        "label": "Пойти сюда",
        "url": "https://absplute-url-to-next-action",
      }
    ]
  }
}
```

`callback` и `button` допускают указание контекста — какой-то дополнительной информации, которая будет возвращена обратно после действия пользователя. Это просто любые свободные данные, прокидываемые как есть (чтобы не хранить лишние состояния на сервере)


### Отправка сообщений по инициативе сервера (обработка очереди)

Бот опрашивает `api.queue.url` с интервалом `api.queue.interval` секунд.
В каждый запрос бот делает методом GET и передедаёт параметр `after=<date_time_utc>` подставляя дату последнего полученного сообщения прошлой итерации:

Например: `https://somehost/messages_queue?after=2024-06-26T09:00:00`

Сервер должен вернуть список сообщений к отправке с датой больше `date_time_utc` в формате:
```json
{
  "status": "complete",
  "response": {
    "items": [
      {
        "datetime": "2024-06-26T10:00:00",
        "tg_chat_id": 12345,
        "data": {"messages": [],"buttons": [],"buttons_header": "","callback": {}}
      }
    ]
  }
}
```

