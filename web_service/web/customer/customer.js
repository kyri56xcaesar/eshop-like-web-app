// Constants
const username = localStorage.getItem("username");
const host = "localhost";

const products_url = "http://"+host+":8081/products/";
const orders_url = "http://"+host+":8082/orders/:"+username;

// Initial method
window.addEventListener("load", (event) => {
    console.log("page is fully loaded");


    const productsElement = document.getElementById("productsDisplay");
    const productsHeader = document.getElementById("products-header");

    const ordersDisplay = document.getElementById("productsDisplay");
    const orders_header = document.getElementById("products-header");

    const basketDisplay = document.getElementById("productsDisplay");
    const basket_header = document.getElementById("products-header");


    basketDisplay.style.display = "none";
    basket_header.style.display = "none";

    ordersDisplay.style.display = "none";
    orders_header.style.display = "none";

    productsElement.style.display = "none";
    productsHeader.style.display = "none";


    loadOrders();
    loadProducts();


});

// default image
const imag = "data:image/jpeg;base64,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";


// Products, orders
var products = [];
var orders = [];
// Instantiate a cart for each user session
var cart = [];

// Display none all divs
const productsElement = document.getElementById("productsDisplay");
productsElement.style.display = "none";
const ordersElement = document.getElementById("ordersDisplay");
ordersElement.style.display = "none";
const productsHeader = document.getElementById("products-header");
productsHeader.style.display = "none";
const ordersHeader = document.getElementById("orders-header");
ordersHeader.style.display = "none";


function loadProducts() {
    // fetch data using axios
    axios.get(products_url)
    .then(function (response) {
        // Handle response
        // console.log(response.data);

        if (response.data != null) {
            products = response.data;
        }

        if (products == null) {
            console.log("Products empty response.");
            return;
        }


        // Display data
        const productsElement = document.getElementById("productsDisplay");

        products.forEach(product => {
            // console.log("Product id retrieved: " + product.id);
            const productDiv = document.createElement('div');
            productDiv.classList.add('product');

            if (product.img == "") {
                product.img = imag;
            }

            productDiv.innerHTML = `
                <div id="product">
                    <h3>${product.title}</h3>
                    <img src="${product.img}" alt="${product.title}">
                    <p>Price: $${product.price}</p>
                    <div class="action-buttons">
                        <button onclick="buyProduct(this.parentElement.parentElement, event)">Buy</button>
                        <input type="number" id="product-amount" name="product-amount" min="1" max="1000" value="1">
                    </div>
                </div>`;

                // Append to the element to the display element
                productsElement.appendChild(productDiv);


                // Also setup cart dictionary.
                cart.push({
                    "id":product.id,
                    "title":product.title,
                    "img":"",// no need to transfer img data
                    "price":product.price,
                    "quantity":0, // this quantity resembles amount user wants to buy
                    "username":product.username
                });
        });
    })
    .catch(function (error) {
        // Handle any errors
        console.log('Error fetching data:', error);
    });


}

function loadOrders() {
    // fetch data using axios
    axios.get(orders_url)
    .then(function (response) {
        // Handle response
        console.log(response.data);
    
        orders = response.data;
    
        if (orders == null) {
            console.log("Orders empty response.");
            return;
        }
    
    
        // Display data
        const ordersElement = document.getElementById("ordersDisplay");
    
        // Populate orders display
        orders.forEach((order, index) => {
            console.log("Order retrieved: " + order);
            const orderDiv = document.createElement('div');
            orderDiv.classList.add('order');
        
        
        
            orderDiv.innerHTML = `
            <div id="order-info">
                <div id="order-customer-info">
                    <div id="order-customer-info-text">
                        <span>User: ${order.username}</span>
                        <span>Price: ${order.total_price}$</span>
                    </div>
                    <div id="order-status">
                        <span>${order.status}</span>
                    </div>
                </div>
                
                <div id="order-products" class="scrollable-container">
                    
                </div>
            </div>`;

            ordersElement.appendChild(orderDiv);

            const order_products = document.getElementById("order-products");
            order_products.classList.add('scrollable-container', 'order-products');
            order_products.id = `order-products-${index}`;
            
            // Populate order products display.
            order.products.forEach(product =>{
                const productDiv = document.createElement('div');
                productDiv.classList.add('order-product');

                if (product.img == "") {
                    product.img = imag;
                }

                productDiv.innerHTML = `
                <div id="order-product">                
                    <div id="order-product-headers">
                        <h4>Title: ${product.title}</h4>
                    </div>
                    <img src="${product.img}" alt="${product.title}">
                    <div id="order-product-footer">
                        <p id="order-product-amount">Amount: ${product.quantity}</p>
                        <p id="order-product-price">Price: ${product.price}$</p>
                    </div>
                </div>`;

                order_products.appendChild(productDiv);

            });



        });
    })
    .catch(function (error) {
        // Handle any errors
        console.log('Error fetching data:', error);
    });
    
}  


// function to display the products
function displayProducts() {
    console.log("Displaying products.")

    const productsDisplay = document.getElementById("productsDisplay");
    const ordersElement = document.getElementById("ordersDisplay");
    const productsHeader = document.getElementById("products-header");
    const ordersHeader = document.getElementById("orders-header");
    const basketHeader = document.getElementById("basket-header");
    const basketDisplay = document.getElementById("basketDisplay");

    basketHeader.style.display = "none";
    basketDisplay.style.display = "none";
    productsHeader.style.display = "flex";
    productsHeader.children[0].innerHTML = "All Products";
    productsDisplay.style.display = "flex";
    ordersHeader.style.display = "none";
    ordersElement.style.display = "none";
}

// function to display the orders
function displayOrders() {
    console.log("Displaying orders.")

    const productsDisplay = document.getElementById("productsDisplay");
    const ordersElement = document.getElementById("ordersDisplay");
    const productsHeader = document.getElementById("products-header");
    const ordersHeader = document.getElementById("orders-header");
    const basketHeader = document.getElementById("basket-header");
    const basketDisplay = document.getElementById("basketDisplay");

    basketHeader.style.display = "none";
    basketDisplay.style.display = "none";
    productsHeader.style.display = "none";
    productsDisplay.style.display = "none";
    ordersHeader.style.display = "block";
    ordersHeader.children[0].innerHTML = "Orders by: " + username;
    ordersElement.style.display = "flex";
}



