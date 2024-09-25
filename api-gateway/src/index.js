const { ApolloServer } = require('apollo-server');
const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: [
      { name: 'user', url: process.env.USER_SERVICE_URL },
      { name: 'string', url: process.env.STRING_SERVICE_URL },
      { name: 'matchup', url: process.env.MATCHUP_SERVICE_URL },
      { name: 'relationship', url: process.env.RELATIONSHIP_SERVICE_URL },
    ],
  }),
});

const server = new ApolloServer({
  gateway,
  subscriptions: false,
});

server.listen({ port: 4000 }).then(({ url }) => {
  console.log(`🚀 Gateway ready at ${url}`);
}).catch(err => {
  console.error('Failed to start the Apollo Gateway:', err);
});
