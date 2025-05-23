#!/bin/sh

HOST=mysql
PORT=3306
USER=root
PASSWORD=arti2002

echo "Ожидаем готовности MySQL на $HOST:$PORT..."

until mysql -h "$HOST" -P "$PORT" -u  root -p"$PASSWORD" --skip-ssl -e "SELECT 1;" 2>&1
do
  echo "MySQL не готова, ждем 1 секунду..."
  sleep 1
done

echo "MySQL готова! Запускаем тесты..."
exec "$@"