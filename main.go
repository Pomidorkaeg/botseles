package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ButtonRequest struct {
	BotToken string `json:"botToken"`
	Command  string `json:"command"`
	Text     string `json:"text"`
}

type BotRequest struct {
	BotToken string `json:"botToken"`
}

func addCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                   // Разрешить доступ с любых доменов
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Разрешить определенные HTTP методы
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Разрешить указанные заголовки
}

func createButton(w http.ResponseWriter, r *http.Request) {
	addCORSHeaders(w) // Добавляем CORS заголовки

	// Обработка OPTIONS запроса
	if r.Method == http.MethodOptions {
		return
	}

	var req ButtonRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bot, err := tgbotapi.NewBotAPI(req.BotToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при создании бота: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Создаем кнопку с заданным текстом
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(req.Text, req.Command),
		),
	)

	// Отправляем сообщение с кнопкой
	msg := tgbotapi.NewMessage(123456789, "Нажмите кнопку ниже:")
	msg.ReplyMarkup = keyboard

	_, err = bot.Send(msg)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при отправке сообщения: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Ответ на запрос
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Кнопка успешно создана!"})
}

func runBot(w http.ResponseWriter, r *http.Request) {
	addCORSHeaders(w) // Добавление CORS заголовков

	// Обработка OPTIONS запроса
	if r.Method == http.MethodOptions {
		return
	}

	var req BotRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Запуск процесса в фоновом режиме
	go func() {
		cmd := exec.Command("go", "run", "bot.go", req.BotToken)
		err := cmd.Start()
		if err != nil {
			log.Printf("Ошибка при запуске бота: %s", err.Error())
		}
	}()

	// Ответ сразу же
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Бот запущен!"})
}

func main() {
	http.HandleFunc("/create_button", createButton)
	http.HandleFunc("/run_bot", runBot)

	// Запуск сервера
	fmt.Println("Запуск сервера на порту 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func botik() {
	// Получаем токен бота
	token := "your_bot_token" // Этот токен передается с клиента

	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	// Логируем, что бот авторизовался
	log.Printf("Авторизован как %s", bot.Self.UserName)

	// Получаем обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	// Обрабатываем обновления
	for update := range updates {
		if update.Message != nil {
			// Отвечаем на сообщение
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я ваш бот.")
			bot.Send(msg)
		}
	}
}
