const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');
const config = require('../../config');
const getServiceUrl = require('../utils/getServiceUrl');
const AuthenticatedDataSource = require('./authenticatedDataSource');

const subgraphs = Object.values(config.SERVICES).map((service) => ({
  name: service.name,
  url: `${getServiceUrl(service.name, service.url)}/graphql`,
}));

// Log the final subgraph URLs before creating the gateway.
console.log('Subgraph service URLs:', subgraphs.map(sg => `${sg.name}: ${sg.url}`));

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({ subgraphs }),
  buildService({ url }) {
    return new AuthenticatedDataSource({ url });
  },
});

module.exports = gateway;
