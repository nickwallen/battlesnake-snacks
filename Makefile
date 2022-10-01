
test:
	go test ./...

run:
	PORT=8000 SNAKE=BATTLE go run ./cmd

battle:
	battlesnake play \
	  -W 11 -H 11 \
	  -g standard \
      --url http://localhost:8000 \
      --url http://localhost:8001 \
      --browser

solo:
	battlesnake play \
	  -W 7 -H 7 \
	  -g standard \
      --url http://localhost:8000 \
      --browser
