const express = require("express");
const cors = require("cors");
const port = process.env.PORT || 5000;


const {connection} = require("./dbConnect.js")


const app = express();
app.use(cors());
app.use(express.json());
app.get("/health", (req, res) => {
  res.send("ok");
});

//get all data
app.get("/products", async (req, res) => {
    try {
        const db = await connection;
        const store = await db.execute(`SELECT * FROM products`);


        res.send(store[0]);
    } catch (error) {
        res.status(500).send(error.message);
    }

});

//get data by id
app.get("/products/:id", async (req, res) => {
   
});

//store new data in a database
app.post("/products", async (req, res) => {

    try {
        const data = req.body;
        const db = await connection;
        const store = await db.execute(`INSERT INTO  products (title, img, price, quantity, username) VALUES (?, ?, ?, ?, ?)`,
            [data.title, data.img, data.price, data.quantity, data.username]
        );
        res.send(store);
    } catch (error) {
        res.status(500).send(error.message);
    }
});

app.listen(port, () => {
    console.log("Products service running in PORT " + port);
});