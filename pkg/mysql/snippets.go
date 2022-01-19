package mysql

import (
	"database/sql"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// snippetModel определяем тип который обертывает пул подключения sql.DB
type SnippetModel struct {
	DB *sql.DB
}

//insert метод для создания новой в базе данных
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	//ниже будет sql запрос, который мы хотим выполнить. мы разделили его на 2 строки
	/// для удобства чтения (пожтому он окружен ``)
	stmt := `INSERT INTO snippets (title, content, created, expires)
VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	//Используя метод exec() из встроенного пула подключений для
	//запроса. Первый параметр это сам sql запрос, за которым следует
	//заголовок заметки, содержимое и строка жизни ззаметки. Этот
	//метод возвращает объект sql.Result, который соделжит некоторые
	//данные о том, что произошло после выполненрия запроса
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	//Используем метод LastInsertID(), чтобы получить последний ID
	//созданной записи из таблицы snippets
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	//Возвращаемый ID имеет тип int64, пожтому мы конвертируем его перед возвратом из метода
	return int(id), nil
}

// get метод для возвращения данных заметок по ее индентификатору
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `select id, title, content, created, expires from snoppets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// используем метод QueryRow() для выполнения SQl запроса,
	//передавая ненадежную переменную id в качестве значения для плейсхолдера(?)
	//возвращается указатель на объект SQL.Row, который содержит данные записи
	row := m.DB.QueryRow(stmt, id)
	//инициализируем указатель на новую струтуру Snippet
	s := &Snippet{}
	//Используйте row.Scan(), чтобы скопировать значения из каждого поля от sql.Row d
	//соответствующее поле в струтуре Snippet. Обратите внимание, что аргумент
	//для row.Scan это указатели на место, куда требуется скопировать данные
	//и количество аргументов должно быть точно таким де, как количество
	//столбцов в таблице базы данных
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		//специиально для жтого случая, мы проверим при помощи функции errors.Is()
		//если запрос был отправлен с ошибкой. Если ошибка обнарудиться, то
		//возвращаем нашу ошибку из модели models.ErrNoRecord
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

//latest метод возвращает 10 наиболее часто используемых заметок
func (m *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
