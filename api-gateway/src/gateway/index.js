const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');
const config = require('../config');
const getServiceUrl = require('../utils/getServiceUrl');

const subgraphs = Object.values(config.SERVICES).map((service) => ({
    name: service.name,
    url: `${getServiceUrl(service.name, service.envUrl)}/graphql`,
}));

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs,
  }),
});

module.exports = gateway;
