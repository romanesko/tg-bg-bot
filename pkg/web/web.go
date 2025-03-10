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
		"actions_url": "/actions",
		"actions_interval": 10

}`

}

func _menuHandler() string {
	return `{
    "messages": [
      {
        "text": "Первое сообщение"
      },
      {
        "text": "Второе сообщение",
        "images": [
          "https://render.fineartamerica.com/images/rendered/default/print/5.5/8/break/images/artworkimages/medium/2/naked-ass-necklace-gene-oryx.jpg", "https://avatars.mds.yandex.net/i?id=f0f808924f403eecf559ccb0e5d74996_l-5228318-images-thumbs&n=13"
        ]
      }
    ],
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

func menuHandler() string {
	return `{

    "buttons_header": "Куда дальше?",
    "buttons": [
    
      {"label": "Шесть", "url": "/subitem"},
	{"label": "Ссылка на внешний ресурс", "link": "https://ya.ru"}
    ]

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

func actionsHandler() string {
	return `{
    "check_user_in_channel": [
      {
        "tg_chat_id": 2042663,
        "tg_channel_name": "@FromCTOtoConsult"
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

		fmt.Printf("TEST API QUERY params: %s, post data: %s\n", params, jsonString)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	}
}
