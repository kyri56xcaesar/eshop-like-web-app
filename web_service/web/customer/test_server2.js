const express = require("express");
const cors = require("cors");




const app = express();
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true}));

app.get("/health", (req, res) => {
    res.send("ok");
})

app.get("/orders/:test", async (req, res) => {

    res.send([
        {
            "id":1,
            "products":[
                {
                    "id":1,
                    "title":"t1",
                    "img":"",
                    "price":1,
                    "quantity":2,
                    "username":"whatever"
                },
                {
                    "id":2,
                    "title":"t2",
                    "img":"",
                    "price":2,
                    "quantity":3,
                    "username":"whatever2"
                },
                {
                    "id":2,
                    "title":"t2",
                    "img":"",
                    "price":2,
                    "quantity":3,
                    "username":"whatever2"
                },
                {
                    "id":2,
                    "title":"t2",
                    "img":"",
                    "price":2,
                    "quantity":3,
                    "username":"whatever2"
                },
                {
                    "id":2,
                    "title":"t2",
                    "img":"",
                    "price":2,
                    "quantity":3,
                    "username":"whatever2"
                },
                {
                    "id":2,
                    "title":"t2",
                    "img":"",
                    "price":2,
                    "quantity":3,
                    "username":"whatever2"
                },
                {
                    "id":2,
                    "title":"t2",
                    "img":"",
                    "price":2,
                    "quantity":3,
                    "username":"whatever2"
                },
                {
                    "id":2,
                    "title":"t2",
                    "img":"",
                    "price":2,
                    "quantity":3,
                    "username":"whatever2"
                },
                {
                    "id":2,
                    "title":"t2",
                    "img":"",
                    "price":2,
                    "quantity":3,
                    "username":"whatever2"
                },
                {
                    "id":2,
                    "title":"t2",
                    "img":"",
                    "price":2,
                    "quantity":3,
                    "username":"whatever2"
                },
                {
                    "id":2,
                    "title":"t2",
                    "img":"",
                    "price":2,
                    "quantity":3,
                    "username":"whatever2"
                },
                {
                    "id":2,
                    "title":"t2",
                    "img":"",
                    "price":2,
                    "quantity":3,
                    "username":"whatever2"
                },
                {
                    "id":2,
                    "title":"t2",
                    "img":"",
                    "price":2,
                    "quantity":3,
                    "username":"whatever2"
                }
                
            ],
            "total_price":344.3,
            "status":"pending",
            "username":"test"
        },
        {
            "id":2,
            "products":[

            ],
            "total_price":6646,
            "status":"Confirmed",
            "username":"test"
        },
        {
            "id":2,
            "products":[
                {
                    "id":4,
                    "title":"t4",
                    "img":"",
                    "price":4,
                    "quantity":4,
                    "username":"whatever4"
                }
            ],
            "total_price":6646,
            "status":"Confirmed",
            "username":"test"
        },
        {
            "id":2,
            "products":[
                {
                    "id":3,
                    "title":"t3",
                    "img":"",
                    "price":3,
                    "quantity":3,
                    "username":"whatever3"
                }
            ],
            "total_price":6646,
            "status":"Confirmed",
            "username":"test"
        }
     
    ]);
});

app.listen("8082", () => {
    console.log("Products service running in PORT 8082");
});
