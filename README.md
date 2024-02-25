# CourseFetcher


Обновление курсов происходит каждые 15 секунд.
Поддерживаются все валюты указанные по адресу:  https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies.json

Для успешного запуска необходимо заполнить пустые поля файла config.yaml в папке config.
В качестве базового hosy и port указаны localhost и 10000 cоответственно.

Так же для успешной работы в папке migration создан файл содержащий ddl query.




1. POST запрос на добавление новой пары валют для отслеживания
POST http://host:port/currency/create

JSON:
{
    "name": "Currency1/Currency2"
}

Пример:
curl -d '{"name": "EUR/JPY"}' -H "Content-Type: application/json" -X POST http://localhost:10000/currency/create

Возвращaет JSON, содержащий уникальный id, по которому можно запросить курс добавленной пары.
HTTP status 201.

JSON:
{
    "id": 30
}

Если попытаться добавить валютную пару, которая уже отслеживается, вернется id созданой пары.

2. GET запрос на получение значение цены валютной пары по уникальному id, полученному при создании или по названию пары.
GET curl -X GET http://host:port/currency/id/{id}
GET curl -X GET http://host:port/currency/code/{currency pair}

Примеры:
curl -X GET http://localhost:10000/currency/id/12
curl -X GET http://localhost:10000/currency/code/EUR/JPY

Возвращает цену и время последнего обновления.

JSON:
{
    "Value":1.08339264,
    "Date":"2024-02-25T23:11:33.464254Z"
}

3. GET запрос для проверки функционирования сервиса
GET POST http://host:port/health

Пример:
curl -X GET http://localhost:10000/health

Возвращает JSON, содержащий bool переменную.

JSON:
{
    "ok": true
}
