
test:
	go test ./...

run:
	PORT=8000 SNAKE=BATTLE go run ./cmd

battle:
	battlesnake play \
	  -W 11 -H 11 \
	  -g standard \
      --name Snake1 --url http://0.0.0.0:8001 \
      --name Snake2 --url http://0.0.0.0:8002 \
      --name Snake3 --url http://0.0.0.0:8003 \
      --output ~/tmp/battlesnake.out \
      --browser

solo:
	battlesnake play \
	  -W 7 -H 7 \
	  -g standard \
      --url http://localhost:8001 \
      --browser
