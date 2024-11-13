require('dotenv').config();
const admin = require('firebase-admin');
const axios = require('axios');
const { ApolloGateway, RemoteGraphQLDataSource } = require('@apollo/gateway');
const { ApolloServer } = require('@apollo/server');
const { startStandaloneServer } = require('@apollo/server/standalone');
const { IntrospectAndCompose } = require('@apollo/gateway');

// Initialize Firebase Admin using environment variable
const serviceAccount = JSON.parse(process.env.FIREBASE_SERVICE_ACCOUNT);
const firebaseConfig = JSON.parse(process.env.FIREBASE_CONFIG);

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount)
});

// Custom DataSource to forward auth headers
class AuthenticatedDataSource extends RemoteGraphQLDataSource {
  willSendRequest({ request, context }) {
    // Forward the auth token to the services
    if (context.token) {
      request.http.headers.set('Authorization', context.token);
    }
  }
}

// Function to exchange custom token for ID token
async function getIdToken(customToken) {
  try {
    const response = await axios.post(
      `https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=${firebaseConfig.apiKey}`,
      {
        token: customToken,
        returnSecureToken: true
      }
    );
    return response.data.idToken;
  } catch (error) {
    console.error('Error exchanging custom token for ID token:', error.response?.data || error.message);
    throw error;
  }
}

// Function to generate custom token and get ID token
async function generateTokens(uid) {
  try {
    const customToken = await admin.auth().createCustomToken(uid);
    const idToken = await getIdToken(customToken);
    console.log(idToken);
    return idToken;
  } catch (error) {
    console.error('Error generating tokens:', error);
    throw error;
  }
}

// Create the gateway instance
const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: [
      { name: 'user-service', url: `${process.env.USER_SERVICE_URL}/graphql` },
      { name: 'relationship-service', url: `${process.env.RELATIONSHIP_SERVICE_URL}/graphql` },
      { name: 'matchup-service', url: `${process.env.MATCHUP_SERVICE_URL}/graphql` },
      { name: 'equipment-service', url: `${process.env.EQUIPMENT_SERVICE_URL}/graphql` },
    ],
  }),
  // Use the custom data source for all services
  buildService({ url }) {
    return new AuthenticatedDataSource({ url });
  },
});

const server = new ApolloServer({
  gateway,
  introspection: process.env.GRAPHQL_PLAYGROUND === 'true',
});

// Start the server and generate tokens
async function startServer() {
  try {
    // Generate token for the specific UID
    const uid = 'FUExVrKHGqTIkNJA2xY64Cvvt6u2';
    const idToken = await generateTokens(uid);

    const { url } = await startStandaloneServer(server, {
      context: async ({ req }) => {
        // If no auth header is present, use the generated token
        const token = req.headers.authorization || `Bearer ${idToken}`;
        return { token };
      },
      listen: { port: parseInt(process.env.PORT || '80') },
    });

    console.log(`
🚀 Server ready at ${url}
📝 Environment: ${process.env.NODE_ENV || 'development'}
📝 Introspection: ${process.env.GRAPHQL_PLAYGROUND === 'true' ? 'enabled' : 'disabled'}
🔗 Connected services:
  - User Service: ${process.env.USER_SERVICE_URL}
  - Relationship Service: ${process.env.RELATIONSHIP_SERVICE_URL}
  - Matchup Service: ${process.env.MATCHUP_SERVICE_URL}
  - Equipment Service: ${process.env.EQUIPMENT_SERVICE_URL}
    `);
  } catch (error) {
    console.error('Failed to start server:', error);
    process.exit(1);
  }
}

startServer();