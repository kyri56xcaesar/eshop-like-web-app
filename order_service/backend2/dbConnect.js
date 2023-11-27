// get the client
const mysql = require('mysql2/promise');


const connect = async ()=>{

    try {
        // create the connection to database
        const connection = await mysql.createConnection({
            // If i use docker, name host the name of the container
            host: 'localhost',
            port: 3308,
            user: 'admin',
            password: 'admin',
            database: 'order_db'
        });

        console.log("Successfully connected to the database.");
        
        const createTable = `
        CREATE TABLE IF NOT EXISTS orders (
            id INT AUTO_INCREMENT PRIMARY KEY,
            products JSON,
            total_price DECIMAL(10, 2),
            status VARCHAR(100),
            username VARCHAR(255)
        );
    `;

    const table = await connection.execute(createTable);

    console.log(table);
        return connection;
    } catch (error) {
        console.log(error);
        throw new Error(error);
    }
    
  
}


const connection = connect();


module.exports = {connection};