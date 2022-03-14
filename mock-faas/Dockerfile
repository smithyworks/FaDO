FROM node:16

WORKDIR /app

COPY package*.json ./

RUN npm ci --only=production

COPY . .

ENV SERVER_NAME=FaDO-mock-faas
ENV SERVER_PORT=9100

EXPOSE 9100

CMD ["node", "index.js"]