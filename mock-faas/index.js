const express = require("express");
const app = express();
app.use(express.json());

const { SERVER_NAME, SERVER_PORT } = process.env;

app.use((req, res) => {
  res.json({
    SERVER_NAME,
    timestamp: new Date().toISOString(),
    method: req.method,
    path: req.path,
    headers: req.rawHeaders,
  });
});

app.listen(SERVER_PORT || 9100, () => {
  console.log(`Server '${SERVER_NAME}' listening on port ${SERVER_PORT || 9100}.`);
});
