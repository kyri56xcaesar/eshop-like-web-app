const { Kafka, Partitioners } = require('kafkajs');


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

module.exports = {
    kafkaProducer: checkOrders
}