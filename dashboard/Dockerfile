FROM node:12.2.0-alpine as dependencies

ADD package.json /tmp/package.json
RUN cd /tmp && yarn install

FROM node:12.2.0-alpine
WORKDIR /usr/src/app

COPY --from=dependencies /tmp/node_modules ./node_modules

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
