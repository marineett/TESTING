#!/bin/bash

# Очищаем кэш npm
npm cache clean --force

# Удаляем node_modules и package-lock.json
rm -rf node_modules
rm -f package-lock.json

# Устанавливаем зависимости заново
npm install 