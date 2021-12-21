FROM node:12.13.1-stretch

RUN npm install @slack/bolt
COPY . .

CMD ["node", "index.js"]

EXPOSE 5000