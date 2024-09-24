require('dotenv').config();
const { ApolloServer } = require('apollo-server');
const { ApolloGateway } = require('@apollo/gateway');
const { services } = require('./config');

const gateway = new ApolloGateway({
  serviceList: services,
});

const server = new ApolloServer({
  gateway,
  subscriptions: false,
});

const port = process.env.PORT || 4000;

server.listen({ port }).then(({ url }) => {
  console.log(`🚀 API Gateway ready at ${url}`);
});