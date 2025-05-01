# Mitter
Mitter is my Twitter-inspired “social network”. 
In early development stage

## Functional

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

## My plans

### Users
- [ ] Following
- [ ] Friends
- [ ] Users' ratings (one user can rate another user only once (rate's editing allowed))

### Mitts
- [ ] Mitts' comments

## How to run

First, clone repo and cd into it

```bash
git clone https://github.com/misshanya/mitter
cd mitter
```

### Docker compose services

- Backend (this app, API)
- PostgreSQL (main DB)
- Redis (right now just for tokens)
- Prometheus 
- Grafana

```bash
docker compose up -d
```

### Build binary
```bash
go build -o server .
```
