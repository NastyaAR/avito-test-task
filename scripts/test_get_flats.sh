#!/bin/bash

# Установить URL и заголовки
URL="http://localhost:80/flat/create"
AUTH_HEADER="Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkX3RpbWUiOjE3MjMxMTQyOTMsInJvbGUiOiJjbGllbnQiLCJ1c2VySUQiOiIwMTkxMzE2Ni02YzVkLTcxODQtYmJkMS1hOWIyODhhY2NlYTcifQ.6EwKY_SaRm3_2M5XK54B5LjWeNonvZVLJBuAxH7retA"

# Количество квартир для создания
TOTAL_FLATS=12000

for ((i=12; i<=TOTAL_FLATS+12; i++)); do
  # Формирование body запроса
  BODY=$(cat <<EOF
{
  "flat_id": $i,
  "house_id": 1,
  "price": 10000,
  "rooms": 4
}
EOF
)

  # Выполнение curl-запроса
  curl -X POST "$URL" \
    -H "$AUTH_HEADER" \
    -H "Content-Type: application/json" \
    -d "$BODY"

  # Вывод информации о статусе
  echo "Created flat with ID: $i"
  
  # Опционально: добавление задержки между запросами, например, 0.1 секунды
  sleep 0.1
done 
