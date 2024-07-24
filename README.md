<a href="https://echo.labstack.com"><img height="80" src="https://cdn.labstack.com/images/echo-logo.svg"></a>

## Backend Golang REST-Api with Clean Architecture

###  Tech Stack
1. [Echo](https://github.com/labstack/echo/)
2. [Gorm](https://github.com/go-gorm/gorm) ( Driver with posgres and mysql )
3. [Validation Input](https://github.com/go-playground/validator)
4. [Swagger](https://github.com/swaggo/echo-swagger) for Documentation, go to [This page](http://localhost:8080/swagger/index.html) to open documentation
5. Encription bcrypt
6. [Viper](https://github.com/spf13/viper) for config
---
### Roadmap
- [ ] Middleware
- [ ] Log
- [ ] Session
- [ ] Rate Limiter / Throttle
- [ ] Translation
- [ ] Queue
- [ ] Shedule
- [ ] Auth with JWT
- [ ] Upload file to AWS s3
---
### Guide
1. Clone this repo
2. Copy config.json
3. Run with command
```
go run cmd/main.go
```