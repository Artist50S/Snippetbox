package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"golangs.org/snippetbox/pkg/mysql"
)

//home page
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) //error 404
		return          //it be must
	}
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	// Используем функцию template.ParseFiles() для чтения файла шаблона.
	// Если возникла ошибка, мы запишем детальное сообщение ошибки и
	// используя функцию http.Error() мы отправим пользователю
	// ответ: 500 Internal Server Error (Внутренняя ошибка на сервере)
	ts, err := template.ParseFiles(files...)
	if err != nil {
		//	app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
	// Затем мы используем метод Execute() для записи содержимого
	// шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	// возможность отправки динамических данных в шаблон.
	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
	//w.Write([]byte("Hello from Snippetbox"))
}
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	//вызываес метод get из модели Snipping для извлечения данных для
	//конкретной зыписи на основк ее id. если подходящей запси не найдено
	//то возвращается 404
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, mysql.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	fmt.Fprintf(w, "%v", s)
}
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	// Используем r.Method для проверки, использует ли запрос метод POST или нет. Обратите внимание,
	// что http.MethodPost является строкой и содержит текст "POST".
	if r.Method != http.MethodPost {
		// Используем метод Header().Set() для добавления заголовка 'Allow: POST' в
		// карту HTTP-заголовков. Первый параметр - название заголовка, а
		// второй параметр - значение заголовка.
		w.Header().Set("Allow", http.MethodPost)
		// Если это не так, то вызывается метод w.WriteHeader() для возвращения статус-кода 405
		// и вызывается метод w.Write() для возвращения тела-ответа с текстом "Метод запрещен".
		// Затем мы завершаем работу функции вызвав "return", чтобы
		// последующий код не выполнялся.
		//w.WriteHeader(405)
		//w.Write([]byte("GET-method isn't allowed"))
		app.clientError(w, http.StatusMethodNotAllowed) // Используем функцию http.Error() для отправки кода состояния 405 с соответствующим сообщением.
		//Error идентиен 2 верхним строчкам
		return
	}
	//Создаем несколько переменных, содержащих тестовые данные.
	title := "the story about the snail"
	content := "snail go,\n eat green"
	expires := "7"
	//передаем данные в метод SnippetModel.Insert(), получая обратно id
	// только что созданной записи в базу данных
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
