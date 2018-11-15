
FROM node:8.11.3

ADD package.json /tmp/package.json
RUN cd /tmp && yarn install
RUN mkdir -p /usr/src/app && cp -a /tmp/node_modules /usr/src/app

# From here we load our application's code in, therefore the previous docker
# "layer" thats been cached will be used if possible
WORKDIR /usr/src/app
ADD . /usr/src/app

RUN npm run build
RUN rm -rf ./build
RUN rm -rf ./test
RUN rm -rf ./src

RUN mv dist/assets/gAPIlogo.PNG dist/assets/gAPIlogo.png
RUN mv public/assets/gAPIlogo.PNG public/assets/gAPIlogo.png

ENV PORT=80

EXPOSE 80

CMD [ "node", "index.js" ]
