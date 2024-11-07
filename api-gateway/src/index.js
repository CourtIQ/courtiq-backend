require('dotenv').config();
const { ApolloGateway } = require('@apollo/gateway');
const { ApolloServer } = require('@apollo/server');
const { startStandaloneServer } = require('@apollo/server/standalone');
const { IntrospectAndCompose } = require('@apollo/gateway');

// Function to format service URLs based on environment
function getServiceUrl(serviceName, port) {
  if (process.env.NODE_ENV === 'development') {
    return `http://${serviceName}:${port}`;
  }
  return process.env[`${serviceName.toUpperCase()}_URL`];
}

// Create the gateway instance
const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: [
      {
        name: 'user-service',
        url: `${process.env.USER_SERVICE_URL}/graphql`
      },
      {
        name: 'relationship-service',
        url: `${process.env.RELATIONSHIP_SERVICE_URL}/graphql`
      },
      {
        name: 'matchup-service',
        url: `${process.env.MATCHUP_SERVICE_URL}/graphql`
      },
      {
        name: 'equipment-service',
        url: `${process.env.EQUIPMENT_SERVICE_URL}/graphql`
      },
    ],
  }),
});

// Initialize the Apollo Server
const server = new ApolloServer({
  gateway,
  introspection: process.env.GRAPHQL_PLAYGROUND === 'true',
});

// Start the server
async function startServer() {
  try {
    const { url } = await startStandaloneServer(server, {
      context: async ({ req }) => ({ token: req.headers.authorization }),
      listen: { 
        port: parseInt(process.env.PORT || '3000')
      },
    });

    console.log(`
🚀 Server ready at ${url}
📝 Environment: ${process.env.NODE_ENV || 'development'}
📝 Introspection: ${process.env.GRAPHQL_PLAYGROUND === 'true' ? 'enabled' : 'disabled'}
🔗 Connected services:
   - User Service: ${getServiceUrl('user_service', 8081)}
   - Relationship Service: ${getServiceUrl('relationship_service', 8082)}
   - Matchup Service: ${getServiceUrl('matchup_service', 8083)}
   - Equipment Service: ${getServiceUrl('equipment_service', 8084)}
    `);
  } catch (error) {
    console.error('Failed to start server:', error);
    process.exit(1);
  }
}

startServer();