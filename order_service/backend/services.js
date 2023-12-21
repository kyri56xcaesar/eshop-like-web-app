//check if products amount is > 0 
// update database if successfull
const { connection } = require("./dbConnect.js");

const handleProducts = async (orders)=> {
    try {
        console.log(orders);
        const db = await connection;

        for await (const object of orders.products) {
            const data = await db.query("SELECT * FROM products WHERE 'id'=?",
                                         [object.product_id]);
            console.log(data);
            const quantity = data[0][0];
            console.log(quantity);
            if (quantity && quantity < amount){
                return false;
            }

            for await (const object of orders.products){
                const data = await db.query("SELECT * FROM products WHERE 'id'=?",
                                         [object.product_id]);
                
                const newQuantity = data[0][0].quantity - object.amount;
                const update = await db.execute("UPDATE products SET quantity = ? WHERE id = ?",
                [newQuantity, obj.product_id]);
            }

            return true;
        }
    } catch (error) {
        console.log(error);
        throw new Error(error);
    }
}

module.exports = { handleProducts }