FROM node:25.6.1-alpine3.23

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

EXPOSE 5173

ENTRYPOINT [ "npm", "run", "dev" ]
