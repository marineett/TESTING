# Используем Node.js как базовый образ
FROM node:18-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем переменные окружения
ENV NODE_ENV=development
ENV CHOKIDAR_USEPOLLING=true
ENV WATCHPACK_POLLING=true
ENV WDS_SOCKET_PORT=3000
ENV NODE_OPTIONS=--openssl-legacy-provider

# Копируем package.json и package-lock.json
COPY package*.json ./

# Устанавливаем зависимости с оптимизациями
RUN npm install

# Копируем остальные файлы проекта
COPY . .

# Запускаем приложение в режиме разработки
CMD ["npm", "start"]