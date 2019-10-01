// Server.js, don't forget to add express & ejs to packages
const express = require("express");

const fs = require("fs");
const app = express();
const port = process.env.PORT || 3003;
const protocol = process.env.FRONTEND_PROTOCOL || "http";
const router = express.Router();

const ApiProtocol = process.env.API_PROTOCOL || "http";
const ApiHost = process.env.API_HOST || "localhost";
const ApiPort = process.env.API_PORT || "8080";
const SocketHost = process.env.SOCKET_HOST || "localhost";
const SocketPort = process.env.SOCKET_PORT || "5000";

const https = require("https");
const http = require("http");

app.use(
  express.static(`${__dirname}/dist`, {
    setHeaders: res => {
      res.setHeader("Api-Base", `${ApiProtocol}://${ApiHost}:${ApiPort}`);

      res.setHeader(
        "Socket-Base",
        `${ApiProtocol}://${SocketHost}:${SocketPort}`
      );
    }
  })
);

app.engine(".html", require("ejs").renderFile);

app.set("views", `${__dirname}/dist`);

router.get("/assets/:file", (req, res) => {
  res.sendFile(`${__dirname}/dist/assets/` + req.params.file);
});

var proxy = require("express-http-proxy");

router.get("/dashboard/*", (req, res) => {
  res.set({
    "Api-Base": `${ApiProtocol}://${ApiHost}:${ApiPort}`,
    "Socket-Base": `${ApiProtocol}://${SocketHost}:${SocketPort}`
  });
  res.sendFile(`${__dirname}/dist/index.html`);
});

app.use("/", router);

app.use("/", proxy(`${ApiProtocol}://${ApiHost}:${ApiPort}`));

switch (protocol) {
  case "https":
    HttpsListen();
    break;
  case "http":
    HttpListen();
    break;
  default:
    HttpListen();
}

function HttpsListen() {
  var privateKey = fs.readFileSync(
    process.env.FRONTEND_CERT_PRIVATE_KEY,
    "utf8"
  );
  var certificate = fs.readFileSync(process.env.FRONTEND_CERT_FILE, "utf8");
  var credentials = {
    key: privateKey,
    cert: certificate
  };

  var httpsServer = https.createServer(credentials, app);
  httpsServer.listen(port);

  console.log(`App running on port ${port} using HTTPS protocol`);
}

function HttpListen() {
  var httpServer = http.createServer(app);
  httpServer.listen(port);

  console.log(`App running on port ${port} using HTTP protocol`);
}