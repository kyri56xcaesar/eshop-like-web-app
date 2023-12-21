const { Kafka, Partitioners } = require('kafkajs');
const { handleProducts } = require('./services');


const kafka = new Kafka({
  clientId: 'products-app',
  brokers: ['localhost:8097'],
  retry: {
    initialRetryTime: 2000,
    retries: 5
  },
});



const producer = kafka.producer({
    allowAutoTopicCreation: true,
    createPartitioner: Partitioners.LegacyPartitioner
});

const checkOrders = async (msg)=> {
    await producer.connect();

    await producer.send({
        topic: 'productsProducer',
        messages: [{
            value: JSON.stringify(msg)
        }]
    });

    await producer.disconnect();
}


const consumer = kafka.consumer({
   groupId: "products-group",
   allowAutoTopicCreation: true,
   
});

const fetchProductsFromOrderTopic = async ()=>{
    try {
        await consumer.connect();
        await consumer.subscribe({topics: ["ordersProducer"]});

        await consumer.run({
            eachMessage: async ({message}) => {
                //console.log(message);
                const jsonMsg = JSON.parse(message.value);
                const result = await handleProducts(jsonMsg);
                if (result) {
                    // return order id and success
                    
                } else {
                    // return order id and reject

                }
            }
        });
    } catch (error) {
        await consumer.disconnect();
        console.log(error.message);
        process.exit(1);
    }
}

setTimeout(async ()=> {
    try {
        await fetchProductsFromOrderTopic();
    } catch (error) {
        console.error(error.message);
        // console.info()
    }
}, 2000);

module.exports = {
    kafkaProducer: checkOrders
}