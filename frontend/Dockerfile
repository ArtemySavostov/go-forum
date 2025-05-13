FROM node:18-alpine AS builder
WORKDIR /forum
COPY package*.json ./
COPY . .
RUN npm install
COPY . .
EXPOSE 3000 
CMD ["npm", "start"]
