prepare_backend:
		cd ${shell pwd}/server && go mod tidy

prepare_frontend:
		cd ${shell pwd}/frontend && pnpm i

prepare:
		make prepare_backend prepare_frontend

start_backend:
		cd ${shell pwd}/server && air &

start_frontend:
		cd ${shell pwd}/frontend && pnpm dev &

start:
		make start_backend start_frontend

stop:
		lsof -t -i:8000 | xargs kill && lsof -t -i:5173 | xargs kill


.PHONY: prepare_backend prepare_frontend prepare start_backend start_frontend start stop
