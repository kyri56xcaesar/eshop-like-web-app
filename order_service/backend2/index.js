const express = require("express");
const cors = require("cors");
const port = process.env.PORT || 5001;


const {connection} = require("./dbConnect.js")


const app = express();
app.use(cors());
app.use(express.json());
app.get("/health", (req, res) => {
  res.send("ok");
});

//get all data
app.get("/orders", async (req, res) => {
    try {
        const db = await connection;
        const store = await db.execute(`SELECT * FROM orders`);


        res.send(store[0]);
    } catch (error) {
        res.status(500).send(error.message);
    }

});

//get data by id
app.get("/orders/:id", async (req, res) => {
  
}); 


const kafka = require('./kafka.js');

//store new data in a database
app.post("/orders", async (req, res) => {

    try {
        // store to database
        const data = req.body;
        const db = await connection;
        const store = await db.execute(`INSERT INTO  orders (products, total_price, status, username) VALUES (?, ?, ?, ?)`,
            [data.products, data.total_price, data.status, data.username]
        );

        // send to kafka
        const msg = {
            id: store[0].insertId,
            products: data.products
        }

        await kafka.kafkaProducer(msg);

        res.send(store);
    } catch (error) {
        res.status(500).send(error.message);
    }
});

app.listen(port, () => {
    console.log("Orders service running in PORT " + port);
});