

function buyProduct(element, e) {
    
    console.log("Buying product, -> Inserting to cart.");

    e.preventDefault();

    // Get the exact product info
    // console.log(element);
    // console.log(element.children);

    if (element.children[3].children[1].value == "") {
        element.children[3].children[1].style.borderColor = "red";

        setTimeout(()=>{
            element.children[3].children[1].style.borderColor = "";

        }, 4000);
        
        return;

    }

    const title = element.children[0].textContent;
    const amount = parseInt(element.children[3].children[1].value);
    
    // console.log(product_info);

    // console.log(title)
    // console.log(amount);


    // Insert them in the cart list.
    cart.find(product => product.title === title).quantity += amount;
    // console.log(cart);

    element.children[3].children[1].value = "";





}


function displayCart() {
    console.log("Displaying cart.");

    const productsDisplay = document.getElementById("productsDisplay");
    const ordersElement = document.getElementById("ordersDisplay");
    const productsHeader = document.getElementById("products-header");
    const ordersHeader = document.getElementById("orders-header");
    const basketHeader = document.getElementById("basket-header");
    const basketDisplay = document.getElementById("basketDisplay");

    
    basketDisplay.style.display = "flex";
    productsHeader.style.display = "none";
    productsDisplay.style.display = "none";
    ordersHeader.style.display = "none";
    ordersElement.style.display = "none";

    basketHeader.style.display = "block";
    basketHeader.children[0].innerHTML = "Your Cart: " + username;

    const basket_info = document.getElementById("basket-info");
    const customer_info = basket_info.children[0].children[0];
    customer_info.children[0].textContent = "User:" + username;
    customer_info.children[1].textContent = "Total: ";

    cleanBasketDisplay();
    updateCart();

}

function updateCart() {

    const basketProducts = document.getElementById("basket-products");

    const basketCheckoutTotalPrice = document.getElementById("basket-customer-info-text").children[1];
    let totalPrice = 0;

    cart.forEach(product => {

        if (product.quantity != 0) {
            
            if (product.img == "") {
                product.img = imag;
            }

            const productDiv = document.createElement('div');
            productDiv.classList.add('product');

            productDiv.innerHTML = `
            <div id="basket-product">                
                <div id="basket-product-headers">
                    <h5>${product.title}</h5>
                    <span class="delete-button" onclick="removeFromBasket(this)">&#10006;</span>
                </div>
                <img src="${product.img}" alt="${product.title}">
                <div id="basket-product-footer">
                    <p id="basket-product-amount">Amount: ${product.quantity}</p>
                    <p id="basket-product-price">Price: ${product.price}$</p>
                </div>
            </div>`;

            basketProducts.appendChild(productDiv);

            totalPrice += product.quantity * product.price;
        }
    

    });

    basketCheckoutTotalPrice.textContent = "Total: " + totalPrice + "$";

    
}

function removeFromBasket(element) {

    // Need to find my title.
    let title = element.parentElement.children[0].textContent;
    console.log(title);


    // update cart list
    cart.find(product => product.title === title).quantity = 0;
    

    let basket = element.parentNode.parentNode.parentNode;

    basket.removeChild(element.parentNode.parentNode);

    // Update price..
    const basketCheckoutTotalPrice = document.getElementById("basket-customer-info-text").children[1];
    let totalPrice = 0;

    cart.forEach(product => {

        if (product.quantity != 0) {
            
            totalPrice += product.quantity * product.price;
        }
    

    });
    basketCheckoutTotalPrice.textContent = "Total: " + totalPrice + "$";
    

}


function cleanBasketDisplay() {

    var basket = document.getElementById("basket-products");
    console.log(basket);

    while (basket.firstChild) {
        basket.removeChild(basket.firstChild);
    }


}


function createOrder(e) {
    console.log("Placing order.");

    e.preventDefault();

    // Get data.
    const basket_customer_info = document.getElementById("basket-customer-info-text");
    const username = basket_customer_info.children[0].textContent.split(":")[1];
    const total_price = parseFloat(basket_customer_info.children[1].textContent.split(" ")[1].split("$")[0]);

    let my_cart = [];
    
    cart.forEach(product => {

        if (product.quantity != 0) {

            product.img = "";
            
            my_cart.push(product);
        }
    

    });

    // console.log(my_cart);
    // console.log(my_cart.length == 0);
    if (my_cart.length == 0) {

        const placeOrderBtn = document.getElementById("place-order-button");
        placeOrderBtn.style.borderColor = "red";

        alert("Cart is empty.");

        setTimeout(() => {
            placeOrderBtn.style.borderColor = "";
        }, 2000);

        return;
    }



    // Setup order
    let order = {
        "products":my_cart,
        "total_price":total_price,
        "status":"Pending",
        "username":username
    };

    console.log(order);

    // Send order request
    let bodyContent = JSON.stringify(order);

    axios.post(orders_url, bodyContent)
        .then((response) => {
            console.log(response.data);

            // let inserted_id = parseInt(response.data.split(":")[1].split(" ")[0]);

        })
        .catch(function (error) {
            console.log('Error creating order:', error);
            return;
        });

    
    
    // Update display
    const ordersElement = document.getElementById("ordersDisplay");

    console.log("ID");

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
    order_products.id = `order-products-${ordersElement.children.length - 1}`;
    
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
                <h5>Title: ${product.title}</h5>
            </div>
            <img src="${product.img}" alt="${product.title}">
            <div id="order-product-footer">
                <p id="order-product-amount">Amount: ${product.quantity}</p>
                <p id="order-product-price">Price: ${product.price}$</p>
            </div>
        </div>`;

        order_products.appendChild(productDiv);

    });



    // Cleanup cart display
    cleanBasketDisplay();



    // Redirect to orders.
    displayOrders();


}


