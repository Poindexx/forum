# Fourm

## Цели

Проект включает создание веб-форума, который позволяет:

- Общение между пользователями.
- Ассоциацию категорий с постами.
- Оценивание постов и комментариев (лайки и дизлайки).
- Фильтрацию постов.
- Сортировка постов.
- Добовление Фото.

## Технологии

- **База данных:** SQLite для хранения данных пользователей, постов, комментариев и т.д.
- **Язык программирования:** Go
- **Аутентификация:** Регистрация и вход пользователей с использованием сессий и куки. Пароли хранятся в зашифрованном виде (bcrypt).
- **Управление контейнерами:** Docker для развертывания проекта.

## Функциональность

- **Регистрация пользователей:** Запрос email, имени пользователя и пароля.
- **Вход пользователей:** Проверка учетных данных и создание сессии.
- **Публикация постов:** Зарегистрированные пользователи могут создавать посты и комментарии.
- **Лайки/дизлайки:** Зарегистрированные пользователи могут оценивать посты и комментарии.
- **Фильтрация:** Фильтрация постов по категориям, созданным постам и лайкам (только для зарегистрированных пользователей).
- **Сортировка:** Сортировка постов по дате, популярности (количеству обсуждений) и категориям.

## Используемые пакеты

- **Стандартные пакеты Go**
- `sqlite3`
- `bcrypt`
- `UUID`

## Objectives (Дополнительные)

- Зарегистрированные пользователи могут создавать посты с изображениями и текстом.
- Просмотр поста должен показывать изображение и текст как пользователям, так и гостям.
- Поддерживаемые форматы изображений: JPEG, PNG, GIF.
- Максимальный размер загружаемых изображений: 20 МБ. При попытке загрузки изображения больше 20 МБ, пользователю отображается сообщение об ошибке.

## Подсказки

- Будьте осторожны с размером изображений.
- Что бы запустить Докер используйте эти команды:
docker image build -f Dockerfile -t forum .
docker container run -p 8080:8080 --detach --name form forum

## Скриншоты

- Основной экран.
[! [подпись] (ссылка_на_изображение)](https://github.com/Poindexx/forum/blob/main/photo/Image%2025.06.2024%20at%2020.10.jpeg)

