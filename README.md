### Route
1. Регистрация нвого пользователя, по умолчанию имеет роль "client"
`http://localhost:8080/api/v1/signup`
1. Авторизация пользователя, по умолчанию существует пользователь login: admin, password: admin с полными правами
`http://localhost:8080/api/v1/signin`
1. Доступен обеим группам
`http://localhost:8080/api/v1/foo`
1. Доступен обеим группам
`http://localhost:8080/api/v1/bar`
1. Доступен только для роли "admin"
`http://localhost:8080/api/v1/sigma`

## Запуск сервера локально: 

1. Клонируем репозиторий 
`git clone https://github.com/SyncretismL/go-auth-rbac`

1. Переходим в корневую папку репозитория 
    `cd go-auth-rbac`

1. Используем утилиту `make` для запуска проекта
- `make run` для запуска проекта


