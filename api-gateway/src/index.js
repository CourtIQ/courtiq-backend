const { ApolloGateway } = require('@apollo/gateway');
const { ApolloServer } = require('apollo-server');

const gateway = new ApolloGateway({
  serviceList: [
    { name: 'user', url: 'http://user-service:8081/query' },
    { name: 'string', url: 'http://string-service:8082/query' },
    { name: 'relationship', url: 'http://relationship-service:8083/query' },
    { name: 'matchup', url: 'http://matchup-service:8084/query' },
  ],
});

const server = new ApolloServer({ gateway });

server.listen().then(({ url }) => {
  console.log(`🚀 Gateway ready at ${url}`);
});
