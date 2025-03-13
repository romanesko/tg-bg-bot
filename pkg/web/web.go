package web

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func configHandler() string {
	return ` {
	"start": "/menu",
	"buttons_header": "Выберите пункт меню ↓",
	"actions": [{"url": "/actions", "interval": 10}],
	"commands" : [{"command": "/start", "description": "Начать с начала"}, {"command":"/balance", "description": "Баланс", "url": "/balance"}],
	"admin_password": "123456"
}`

}

func menuHandler() string {
	return `{
    "messages": [
      {
        "text": "Первое сообщение"
      },
      {
        "text": "Второе сообщение с картинкой",
        "images": [
          "https://render.fineartamerica.com/images/rendered/default/print/5.5/8/break/images/artworkimages/medium/2/naked-ass-necklace-gene-oryx.jpg", "https://avatars.mds.yandex.net/i?id=f0f808924f403eecf559ccb0e5d74996_l-5228318-images-thumbs&n=13"
        ]
      },
      {
		"text": "Ссылка без превью: https://ya.ru"
		},
 {
		"text": "Ссылка с превью: https://ya.ru",
		"show_preview": true
		}
    ],
    "buttons_header": "Куда дальше?",
    "buttons": [
      {"label": "Один", "url": "/subitem"}, 
      {"label": "Два", "url": "/subitem2"}, 
	  {"label": "Три", "url": "/relative-url-to-next-action"},
      {"label": "Четыре", "url": "https://absplute-url-to-next-action"},
	  {"label": "Пять", "url": "/relative-url-to-next-action"},
      {"label": "Шесть", "url": "https://absplute-url-to-next-action"},
      {"label": "Ссылка на внешний ресурс", "link": "https://ya.ru"}
    ],
	"buttons_rows": [2,5]

  }`
}

func _menuHandler() string {
	return `{
	"messages": [
      {
        "text": "https://avatars.mds.yandex.net/i?id=f0f808924f403eecf559ccb0e5d74996_l-5228318-images-thumbs&n=13"
      }],
    "buttons_header": "Куда дальше?",
    "buttons": [
      {"label": "Один", "url": "/subitem"}, 
      {"label": "Два", "url": "https://absplute-url-to-next-action"},
	  {"label": "Три", "url": "/relative-url-to-next-action"},
      {"label": "Четыре", "url": "https://absplute-url-to-next-action"},
	  {"label": "Пять", "url": "/relative-url-to-next-action"},
      {"label": "Шесть", "url": "https://absplute-url-to-next-action"},
      {"label": "Ссылка на внешний ресурс", "link": "https://ya.ru"}
    ],
	"buttons_rows": [2,5]
  }`
}

func subItem() string {
	return `{
        "messages": [
            {
                "text": "Подпункт"
            }
        ],
        "buttons_header": "Куда дальше?",
        "buttons": [
            {
                "label": "В начало",
                "url": "/menu"
            }
        ],
        
        "callback": {
            "url": "/calback-url",
            "context": {}
        }
}
    `
}

func subItem2() string {
	return `{
        "messages": [
            {
                "text": "<b>Детали бодиграфа<\/b>\n\n1981-04-09 08:45:00, Москва, Россия\n\nТип: Проектор\n\nПрофиль: 5\\1 - Еретик\\Исследователь\n\nАвторитет: Селезеночный\n\nОпределенность: Цельная\n\nИнкарнационный крест: Горна (Глашатая) 1 ( 51\/57 | 61\/62 )\n\nПеременные: PLL-DLRКартинка: https:\/\/bodygraph.online\/bodygraphs\/67cf78c7dacd20.45447969.png",
                "images": [
                    "https://bodygraph.online/bodygraphs/67cf78c7dacd20.45447969.png"
                ]
            }
        ]
    }`
}

func actionsHandler() string {
	return `{
    "check_user_in_channel": [
      {
        "tg_chat_id": 2042663,
        "tg_channel_name": "@FromCTOtoConsult"
      },
      {
        "tg_chat_id": 217826967,
        "tg_channel_name": "@FromCTOtoConsult"
      },
	{
        "tg_chat_id": 217826967,
        "tg_channel_name": "-1002391915441"
      }
    ]
  }`
}

func balance() string {
	return `{
        "messages": [
            {
                "text": "Ваш баланс: 1000 рублей"
            }
        ]
    }`
}

func Init() {
	// Регистрируем обработчик для маршрута "/"
	http.HandleFunc("/config", handlerWrapper(configHandler))
	http.HandleFunc("/menu", handlerWrapper(menuHandler))
	http.HandleFunc("/actions", handlerWrapper(actionsHandler))
	http.HandleFunc("/subitem", handlerWrapper(subItem))
	http.HandleFunc("/subitem2", handlerWrapper(subItem2))
	http.HandleFunc("/balance", handlerWrapper(balance))

	// Запуск сервера на порту 8080
	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed to start:", err)
	}
}

func handlerWrapper(fun func() string) http.HandlerFunc {
	var data = fun()
	var resp = fmt.Sprintf(`{
            "status": "complete",
            "response": %s }`, data)
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		// Read the body of the request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Failed to read request body")
		}
		// It's good practice to close the body reader
		defer r.Body.Close()

		// Convert the body to a string and print it
		jsonString := string(body)

		log.Printf("TEST WEB API INCOMING REQUEST :: params: %s, post data: %s\n", params, jsonString)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	}
}
