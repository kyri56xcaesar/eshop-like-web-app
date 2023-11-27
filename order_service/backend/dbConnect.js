// get the client
const mysql = require('mysql2/promise');


const connect = async ()=>{

    try {
        // create the connection to database
        const connection = await mysql.createConnection({
            // If i use docker, name host the name of the container
            host: 'localhost',
            port: 3310,
            user: 'admin',
            password: 'admin',
            database: 'products_db'
        });

        console.log("Successfully connected to the database.");

        const createTable = `
            CREATE TABLE IF NOT EXISTS products (
                id INT AUTO_INCREMENT PRIMARY KEY,
                title VARCHAR(255),
                img VARCHAR(255),
                price DECIMAL(10, 2),
                quantity INT,
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