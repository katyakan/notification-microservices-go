FROM node:20-alpine

WORKDIR /usr/src/app

COPY package*.json ./
COPY tsconfig*.json ./
COPY .env .env
COPY nest-cli.json ./

RUN npm install

COPY apps/notification-service ./apps/notification-service
COPY libs ./libs

RUN npm run build notification-service

CMD ["node", "dist/apps/notification-service/main"]