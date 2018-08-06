// Server.js, don't forget to add express & ejs to packages
const express = require("express");

const fs = require("fs");
const app = express();
const port = process.env.PORT || 3003;
const protocol = process.env.FRONTEND_PROTOCOL || "http";
const router = express.Router();

const https = require("https");
const http = require("http");

app.use(express.static(`${__dirname}/dist`)); // set the static files location for the static html

app.engine(".html", require("ejs").renderFile);

app.set("views", `${__dirname}/dist`);

router.get("/assets/:file", (req, res) => {
  res.sendFile(`${__dirname}/dist/assets/` + req.params.file);
});

router.get("/*", (req, res) => {
  res.sendFile(`${__dirname}/dist/index.html`);
});

app.use("/", router);

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
  var credentials = { key: privateKey, cert: certificate };

  var httpsServer = https.createServer(credentials, app);
  httpsServer.listen(port);

  console.log(`App running on port ${port} using HTTPS protocol`);
}

function HttpListen() {
  var httpServer = http.createServer(app);
  httpServer.listen(port);

  console.log(`App running on port ${port} using HTTP protocol`);
}
