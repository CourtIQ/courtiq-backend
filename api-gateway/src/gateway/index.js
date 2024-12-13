const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');
const config = require('../config');
const getServiceUrl = require('../utils/getServiceUrl');
const { RemoteGraphQLDataSource } = require('@apollo/gateway');

const subgraphs = Object.values(config.SERVICES).map(service => ({
  name: service.name,
  url: `${getServiceUrl(service.name, service.url)}/graphql`,
}));

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({ subgraphs }),
  buildService({ name, url }) {
    return new RemoteGraphQLDataSource({
      url,
      willSendRequest({ request, context }) {
        // This is your chance to forward claims to the downstream service
        if (context.user) {
          request.http.headers.set('X-User-Claims', JSON.stringify(context.user));
        }
      },
    });
  }
});

module.exports = gateway;
