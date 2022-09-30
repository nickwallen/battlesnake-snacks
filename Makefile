
run:
	go run ./cmd

battle:
	battlesnake play \
	  -W 11 -H 11 \
	  -g standard \
      --url http://localhost:8000 \
      --url http://localhost:8001 \
      --browser