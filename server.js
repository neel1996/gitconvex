const globalAPI = require("./global/globalAPIHandler");
const express = require("express");
const path = require("path");

const dotenv = require("dotenv").config();

const app = globalAPI;

app.use(express.static(path.join(__dirname, "build")));

app.get("/*", (req, res) => {
  res.sendFile(path.join(__dirname, "build", "index.html"));
});

globalAPI.listen(process.env.GITCONVEX_PORT || 9001, (err) => {
  if (err) {
    console.log(err);
  }
  console.log("GitConvex API connected!");
});
