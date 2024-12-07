// index.js
require('dotenv').config();
const { ApolloServer } = require('@apollo/server');
const { startStandaloneServer } = require('@apollo/server/standalone');
const config = require('./config');
const gateway = require('./gateway');
const authMiddleware = require('./middleware/auth'); // Ensure this file exports a function like the one above

async function startServer() {
  try {
    const server = new ApolloServer({
      gateway,
      introspection: config.GRAPHQL_PLAYGROUND,
    });

    const { url } = await startStandaloneServer(server, {
      // Use the authMiddleware as your context function:
      context: authMiddleware,
      listen: { port: config.PORT },
    });

    console.log(`🚀 Server ready at ${url}`);
  } catch (error) {
    console.error('Failed to start server:', error.message);
    process.exit(1);
  }
}

startServer();