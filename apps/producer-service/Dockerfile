# Stage 1 — builder
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
COPY nest-cli.json ./
COPY tsconfig*.json ./
COPY .env .env
COPY apps ./apps
COPY libs ./libs

COPY apps/producer-service/src/main.ts ./apps/producer-service/

RUN npm install
RUN npm run build:producer-service

# Stage 2 — runtime
FROM node:20-alpine

WORKDIR /app

COPY --from=builder /app/dist /app/dist
COPY --from=builder /app/node_modules /app/node_modules
COPY package*.json ./

CMD ["node", "dist/apps/producer-service/main.js"]