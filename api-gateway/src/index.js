require('dotenv').config();
const { admin } = require('./middleware/auth'); // Ensure Firebase is initialized
const { ApolloServer } = require('@apollo/server');
const { startStandaloneServer } = require('@apollo/server/standalone');
const config = require('./config');
const gateway = require('./gateway');
const { generateIdToken } = require('./utils/firebaseUtils');

// Main function to start the server
async function startServer() {
  try {
    const server = new ApolloServer({
      gateway,
      introspection: config.GRAPHQL_PLAYGROUND,
    });

    // Example: Generate an ID token for a test user (for debugging)
    const uid = 'test-uid'; // Replace with a valid UID from your Firebase Auth
    const idToken = await generateIdToken(uid);
    console.log(`Generated ID Token for testing: ${idToken}`);

    const { url } = await startStandaloneServer(server, {
      context: async ({ req }) => {
        const token = req.headers.authorization || null;
        return { token }; // Pass token to resolvers
      },
      listen: { port: config.PORT },
    });

    console.log(`🚀 Server ready at ${url}`);
  } catch (error) {
    console.error('Failed to start server:', error.message);
    process.exit(1);
  }
}

startServer();
