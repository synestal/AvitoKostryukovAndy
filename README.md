#  Основные функции и заметки по запуску:
  	Docker:	docker-compose build;					docker-compose up -d
  	Тест на прием (интеграционный), 				запуск из корня: cd app/test; go run initTests.go; cd..; cd..
  	Линтер (параметры находятся в app\config), 			запуск из корня: cd app/cmd; go run analyser.go; cd..; cd..
   
  	Для запуска из-под localhost заменить yaml файл с конфиграцией, прописав в ней localhost
# Стек
  Go 1.22, Redis-alpine, potgreSQL-16, Docker, Goland-IDE
#  Команды внутри проекта
  	GET      получение баннера                  http://localhost:8080/user_banner?tag_id=8&feature_id=15&use_last_revision=true&admin_token=25
	POST     внести пользователя                http://localhost:8080/set_admin?id=25&state=true
	POST     создать баннер                     http://localhost:8080/banner?admin_token=10&feature_id=15&tag_ids=22,12&content=notebooklovers,simpledescr,http://aboba.com&is_active=true
	PATCH    изменить баннер                    http://localhost:8080/banner/{id}?admin_token=25&feature_id=100&tag_ids=100,101&content=avitolovers,descr,http://avito.com&is_active=true&id=3
	DELETE   удалить баннер                     http://localhost:8080/banner/{id}?admin_token=10&id=3
	GET      получить баннер с фильтром         http://localhost:8080/banner?admin_token=25&feature_id=15&tag_id=15&content=5&offset=0
	DELETE   удалить отложенно                  http://localhost:8080/delete?admin_token=25&feature_id=15
	GET      получить историю                   http://localhost:8080/history?admin_token=25&id=15
	PATCH    изменить из истории                http://localhost:8080/history?admin_token=25&id=15&number=8
#  Уточнения
##  Кеширование.
Кеширование реализованно на Redis - запрос кеширует в бд только раз в 5 минут индивидуально для каждого запроса в теле самого обработчика, а не фоном раз в 5 минут кешируя всю бд (если запроса нет час, то обновления не будет пока не будет вызван запрос).
##  Фичи и теги.
  Фичи и теги типа integer и integer[] соответственно
##  Дополнительный функционал.
  Отложенное удаление, получение истории, авторизация пользователя, на их основе изменен api.yaml
#  Архитектура
##  Начнем описывать архитетуру с построения DFD диаграммы потоков информации с декомпозицией. DFD в данном случае лучше IDF или других стандартов, так как выгодно описывает взаимодействие сущностей в потоке.
##  Первый уровень
![DFD первый уровень декомпозиции](https://github.com/synestal/AvitoKostryukovAndy/blob/main/%D0%90%D1%80%D1%85%D0%B8%D1%82%D0%B5%D0%BA%D1%82%D1%83%D1%80%D0%B0/DFD%20%D0%B4%D0%B8%D0%B0%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D0%B0%20%D0%BF%D0%BE%D1%82%D0%BE%D0%BA%D0%BE%D0%B2%2C%201-%D1%8B%D0%B9%20%D1%83%D1%80%D0%BE%D0%B2%D0%B5%D0%BD%D1%8C%20%D0%B4%D0%B5%D0%BA%D0%BE%D0%BC%D0%BF%D0%BE%D0%B7%D0%B8%D1%86%D0%B8%D0%B8%20%D0%B2%20Ramus.jpg)

##  И второй уровень
![DFD второй уровень декомпозиции](https://github.com/synestal/AvitoKostryukovAndy/blob/main/%D0%90%D1%80%D1%85%D0%B8%D1%82%D0%B5%D0%BA%D1%82%D1%83%D1%80%D0%B0/DFD%20%D0%B4%D0%B8%D0%B0%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D0%B0%20%D0%BF%D0%BE%D1%82%D0%BE%D0%BA%D0%BE%D0%B2%2C%202-%D0%BE%D0%B9%20%D1%83%D1%80%D0%BE%D0%B2%D0%B5%D0%BD%D1%8C%20%D0%B4%D0%B5%D0%BA%D0%BE%D0%BC%D0%BF%D0%BE%D0%B7%D0%B8%D1%86%D0%B8%D0%B8%20%D0%B2%20Ramus.jpg)

Таким образом, можно увидеть три слоя работы программы - слой handler, слой functions и слой DAO (взаимодействия с бд).

##  Перейдем к детальному рассмотрению работы программы.
Из источника (например, браузера) приходит запрос на выбранный сокет, где работает программа. Далее, в зависимости от типа запроса и его заголовка, сам запрос пересылается в handler (обработчик запроса, слой handler). В сущности handler происходит изъятие информации из http запроса и необходимая проверка на целостность данных для конкретного типа водных данных (например, поле features равно слову, а не числу). На этом этапе может быть сгенерирован ответ об ошибке.
Далее обработчк вызывает сущности обработки запроса и ожидает ответа от нее. Затем функция обработки (слой functions) производит операции для удволетворения http-запроса, например, проверяет является человек админом и показывает баннер. В этом слое сообщаются ошибки и отправляются в handler как отдельный объект. Обращаясь к функции показа баннера, программа переходит к новому слою - слою DAO.
В слое DAO происходит общение с бд postgres. Выполняя все нужные запросы, слой DAO генерирует ответ и ошибку, которая потом передается слою functions, который передает ее header слою. Сам handler слой формирует ответ на основе error и answer из слоя functions конечному пользователю.

##  Опишем схему бд postgres AvitoDB.

Изначально была выбрана модель разбития информации о баннере и о его фичах в две разные сущности.
![Схема бд, первая ревизия](https://github.com/synestal/AvitoKostryukovAndy/blob/main/%D0%90%D1%80%D1%85%D0%B8%D1%82%D0%B5%D0%BA%D1%82%D1%83%D1%80%D0%B0/%D0%A1%D1%85%D0%B5%D0%BC%D0%B0%20%D0%91%D0%94%20postres%20%D0%B2%20Erwin.jpg)

Но она проиграла второму решению по скорости работы и простоте разработки кода, а так же, отсутствия необходимости решать проблемы связи многие-ко многим на физическом уровне БД, поэтому в проекте имплементируется именно второй подход.
![Схема бд, вторая ревизия](https://github.com/synestal/AvitoKostryukovAndy/blob/main/%D0%90%D1%80%D1%85%D0%B8%D1%82%D0%B5%D0%BA%D1%82%D1%83%D1%80%D0%B0/%D0%90%D1%80%D1%85%D0%B8%D1%82%D0%B5%D0%BA%D1%82%D1%83%D1%80%D0%B0%20postgres%20%D1%84%D0%B8%D0%BD%D0%B0%D0%BB%D1%8C%D0%BD%D0%B0%D1%8F%20Erwin.jpg)

Предполагается, что уже существует некая таблица, содержащая данные о пользователях и в нее заносит и изменяет данные другой сервис, но для простоты отладки и работы, был написан способ добавления нового пользователя с флагом - таблица user_tokens.
В таблице banners_storage содержится вся информация о баннере и его фичах, тегах, времени обновления и прочем. В таблице history_banenrs содержится информация об истории баннера - 3 последних обновления и текущее состояние. В таблице delayed_deletions реализовано отложенное удаление

#  Дополнительные задания.
##	1. Были рассмотрены 2 различные архитектуры, проведены тесты, изменен метод кеширования на индивидуальное кеширование. 
##	2. Проведено нагрузочное тестирование в JMeter.
  Изначально был выбран Postman, но он не поддерживал заданное количество транзакций, поэтому был выбран Jmeter чтобы смоделировать поведение системы, сгенерируем таблицу в 1000 ячеек. Далее начнем тест с 1000 RPS. Большинство запросов должны быть GET от обычнвх пользователей, так как, это баннеры на авито, значит, на 100 GET будет 3 POST, 3 PATCH 3 DELETE и еще 10 GET на историю изменений. Обозначим персентиль - 33/1/1/1/3. В JMeter выставим нагрузку в 8 одновременных потоков.
##	3. Добавлена возможность получения и изменения баннеров на основе их истории изменения. 
В таблице заведен триггер, который при обновлении или создании записей о баннерах добавит новую запись в таблицу с историей и если нужно удалит самую старую строку (5-ую по счету). При изменении на основе истории нужен индекс того что изменяем и под каким номер он идет в выводе. Тогда происходит удаление соответсвующей строки из истории и вставка элемента в бащовую таблицу. Таким образом, происходит изменение без потери информации.
##	4. Добавлена возможность отложенного удаления по фиче или тегу. Записи отправляются в отдельную таблицу, где они по триггеру добавления удаляются.
##	6. Линтер добавлен в формате go файла с конфигурацией в .golangci.yml
