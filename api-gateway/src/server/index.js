const { ApolloServer } = require('@apollo/server');
const gateway = require('../gateway');
const config = require('../config');

const server = new ApolloServer({
  gateway,
  introspection: config.GRAPHQL_PLAYGROUND,
});

module.exports = server;
