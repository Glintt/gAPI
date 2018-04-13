//Server.js, don't forget to add express & ejs to packages
const express = require('express')

const app = express()
const port = process.env.PORT || 3003
const router = express.Router()

const API_HOST = process.env.API_HOST || "localhost";
const API_PORT = process.env.API_PORT || "8080";
const SOCKET_HOST = process.env.SOCKET_HOST || "localhost";
const SOCKET_PORT = process.env.SOCKET_PORT || "5000";
const cookieMaxAge = 60*60;

app.use(express.static(`${__dirname}/dist`)) // set the static files location for the static html

app.engine('.html', require('ejs').renderFile)

app.set('views', `${__dirname}/dist`)

router.get('/assets/:file', (req, res, next) => {
  res.sendFile(`${__dirname}/dist/assets/` + req.params.file)
})

router.get('/*', (req, res, next) => {
  res.cookie('API_HOST', API_HOST, { maxAge: cookieMaxAge, httpOnly: false });
  res.cookie('API_PORT', API_PORT, { maxAge: cookieMaxAge, httpOnly: false });
  res.cookie('SOCKET_HOST', SOCKET_HOST, { maxAge: cookieMaxAge, httpOnly: false });
  res.cookie('SOCKET_PORT', SOCKET_PORT, { maxAge: cookieMaxAge, httpOnly: false });
  res.sendFile(`${__dirname}/dist/index.html`);  
})

app.use('/', router)

app.listen(port)
console.log('App running on port', port)