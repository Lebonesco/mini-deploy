FROM node:8

ARG repo
ARG name
ARG port

# Downloads user repo
RUN git clone $repo

WORKDIR /$name
COPY package*.json ./

RUN npm install

EXPOSE $port

CMD ["npm", "start"]