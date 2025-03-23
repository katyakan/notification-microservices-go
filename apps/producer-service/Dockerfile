FROM node:18-alpine

WORKDIR /usr/src/app

COPY package*.json ./
COPY tsconfig*.json ./
COPY nest-cli.json ./

RUN npm install

COPY apps/producer-service ./apps/producer-service
COPY libs ./libs

RUN npm run build producer-service

CMD ["node", "dist/apps/producer-service/main"]