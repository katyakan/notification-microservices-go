FROM node:20-alpine

WORKDIR /usr/src/app

COPY package*.json ./
COPY tsconfig*.json ./
COPY .env .env
COPY nest-cli.json ./

RUN npm install

COPY apps/consumer-service ./apps/consumer-service
COPY libs ./libs

RUN npm run build consumer-service

CMD ["node", "dist/apps/consumer-service/main"]