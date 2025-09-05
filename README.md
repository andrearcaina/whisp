# whisp

<img src="static/favicon.ico" width="160" />

Real time simple anonymous chat app where anyone can send messages without the need to sign up or log in

### tech stack
- Backend: Go, Gin, Gorilla WebSockets
- Frontend: Templ, Alpine.js, TailwindCSS
- Database: NeonDB (PostgreSQL)
    - Goose (DB Migrations) 
    - SQLc (DB interactions)
    - pgx (DB persistence)
- Hosting & Infrastructure: GCP, Cloud Run, Docker
- CI/CD Tools: Makefile, GitHub Actions, Air (Go Live Reload)

### todo (in no particular order)
- [ ] add auth and login maybe (but still keep anonymity)
- [ ] allow users to send images
- [ ] allow users to send gifs (using Tenor API)
