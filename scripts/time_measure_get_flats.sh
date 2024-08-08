#!/bin/bash

# Установим переменные
URL="http://localhost:8081/house/1"
AUTH_HEADER="Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVkX3RpbWUiOjE3MjMwMjc5MjUsInJvbGUiOiJtb2RlcmF0b3IiLCJ1c2VySUQiOiIwMTkxMmM0MC1iYmUyLTc2M2UtYTAwZi1lNTkyYmY3NGNhMTcifQ.IIm2LylqFUzUr1KiXe2UPr9lct_zk6VLsEHhOjYFmZs"
TOTAL_REQUESTS=1000
TOTAL_TIME=0

# Выполнение запросов
for ((i=1; i<=TOTAL_REQUESTS; i++)); do
  # Замер времени выполнения запроса
  START_TIME=$(date +%s%N)  # Начало замера времени в наносекундах
  curl -s -X GET "$URL" -H "$AUTH_HEADER" > /dev/null  # Выполнение запроса
  END_TIME=$(date +%s%N)    # Конец замера времени в наносекундах

  # Расчет продолжительности в миллисекундах
  DURATION=$((($END_TIME - $START_TIME) / 1000000))  # Время в миллисекундах
  TOTAL_TIME=$((TOTAL_TIME + DURATION))  # Сумма времени
done

# Вычисление среднего времени
AVERAGE_TIME=$((TOTAL_TIME / TOTAL_REQUESTS))

# Вывод среднего времени
echo "Average time for $TOTAL_REQUESTS requests: $AVERAGE_TIME ms" 
