# Mitter
[![wakatime](https://wakatime.com/badge/user/6c2e820c-673b-4690-9190-7b15c368b37f/project/a0a1543a-1b5c-4206-814b-c661a923cec8.svg?style=for-the-badge)](#)

Mitter is my Twitter-inspired “social network”. 
In early development stage

## Functionality

### Users
- Sign-Up
- Sign-In
- Update profile (change name)
- Change password
- Delete account

### Mitts
- Create mitt
- Update mitt (change content)
- Get user's mitts
- Get mitt by id
- Like mitt
- Delete mitt

## API Documentation

Swagger docs are located in `docs/swagger.json` or `docs/swagger.yaml`

You can use interactive swagger on `/swagger/index.html` (if my server is running on localhost:8080, the swagger url will be `http://localhost:8080/swagger/index.html`)

## Stack

- Go
- PostgreSQL as main DB
- Redis (right now just for auth tokens)
- Prometheus
- Grafana

## How to run

First, clone repo and cd into it

```bash
git clone https://github.com/misshanya/mitter
cd mitter
```

### Docker compose
```bash
docker compose up -d
```

### Build binary
```bash
go build -o server .
```

## My plans

### Users
- [ ] Following
- [ ] Friends
- [ ] Users' ratings (one user can rate another user only once (rate editing allowed))

### Mitts
- [ ] Mitts' comments

