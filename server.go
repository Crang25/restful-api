package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Book model
type Book struct {
	// Id книги
	ID string `json:"id"`
	// Названик книги
	Title string `json:"title"`
	// Автор книги
	Author *Author `json:"author"`
}

// Author model
type Author struct {
	// Фамилия автора
	FirstName string `json:"firstName"`
	// Имя автора
	LastName string `json:"lastName"`
}

// Срез с книгами
var books []Book

// Отправляет в JSON формате список всех книг. Пользователь делает GET запрос по url your_domain/books
func getBooks(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type, определяющим с каким типом данных будем работать, а именно - json
	w.Header().Set("Content-Type", "application/json")
	// Кодируем срез книг books в JSON и отправляем пользователю
	_ = json.NewEncoder(w).Encode(books)
}

// Отправляет в JSON формате книгу с id, заданным в url. Пользователь делает GET запрос по url your_domain/books/id
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Передаем в ф-ю mux.Vars() наш запрос r, которая возвращает распарсенный url, который отправил пользователь
	// Например, получение книги с id 1: your_domain/books/1. Vars(r)["id"] вернет - 1
	params := mux.Vars(r)
	for _, item := range books {
		if item.ID == params["id"] {
			// Формируем ответ в json и отправляем пользователю
			_ = json.NewEncoder(w).Encode(item)
			return
		}
	}
	// Если нужная книга не найдена, отправляем пользователю json пустой книги
	_ = json.NewEncoder(w).Encode(&Book{})
}

// Создает новую книгу и добавляет в books. Пользователь делает POST запрос в JSON формате и отправляет по адресу
// your_domain/books. Сервер генерирует ID для новой книги и добавляет книгу в список книг.
// Дальше отправляет полный JSON созданной книги пользователю
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Инициализация объекта структуры Book
	var book Book
	// Декодируем полученную книгу от пользователя в JSON формате в структуру Book и записываем в ранее объявленную переменную
	_ = json.NewDecoder(r.Body).Decode(&book)
	// Генерируем рандомный ID для новой книги
	book.ID = fmt.Sprint(rand.Intn(1000000))
	// Добавляем созданную пользователем книгу в наш срез книг
	books = append(books, book)
	// Формируем JSON ответ в виде созданной книги и отправляем пользователю
	_ = json.NewEncoder(w).Encode(book)
}

// Пользователь делает PUT запрос на сервер по адресу your_domain/books/id, где id - id книги, парамеры которой
// пользователь хочет изменить, и отправляет JSON книги с уже измененными параметрами
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			// Удаляем книгу с выбранным пользователем id
			books = append(books[:index], books[index+1:]...)
			var book Book
			// Записываем обновленную пользователем книгу в список книг
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			// Отправляем JSON обновленной книги пользователю
			_ = json.NewEncoder(w).Encode(book)
			return
		}
	}
	// Если нужная книга не найдена, отправляем json со списком всех книг books
	_ = json.NewEncoder(w).Encode(books)
}

// Пользователь делает DEL запрос по адресу your_domain/books/id, где id - id удаляемой книги
// А сервер просто удаляет книгу, с выбранным id из списка книг books
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			// Удаляем книгу с нужным id из списка книг books
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	// Если нужная книга не найдена, отпавляем пользователю json со списком всех книг books
	_ = json.NewEncoder(w).Encode(books)
}

func route() {
	// Порт, на котором будет работать сервер
	port := ":8000"

	// Объявляем новый роутер
	router := mux.NewRouter()
	// Устанавливаем обработчики
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	log.Fatalln(http.ListenAndServe(port, router))
}

func main() {
	// При каждом запуске программы будут генерироваться разные числа
	rand.Seed(time.Now().UnixNano())

	// Добавляем книги для теста
	books = append(books, Book{
		"1",
		"Война и мир",
		&Author{
			"Толстой",
			"Лев"}},
		Book{
			"2",
			"Отцы и дети",
			&Author{
				"Тургенев",
				"Иван"}})

	// Запуск сервера
	route()

}
