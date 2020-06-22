const globalAPI = require("./global/globalAPIHandler");

globalAPI.listen(9001, (err) => {
  if (err) {
    console.log(err);
  }
  console.log("GitConvex API connected!");
});
