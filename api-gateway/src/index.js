// index.js
require('dotenv').config();
const { ApolloServer } = require('@apollo/server');
const { startStandaloneServer } = require('@apollo/server/standalone');
const config = require('./config');
const gateway = require('./gateway');
const authMiddleware = require('./middleware/authMiddleware');

async function startServer() {
  const server = new ApolloServer({
    gateway,
    introspection: config.GRAPHQL_PLAYGROUND,
  });

  const { url } = await startStandaloneServer(server, {
    context: authMiddleware,
    listen: { port: config.PORT },
  });

  console.log(`🚀 Server ready at ${url}`);
}

startServer().catch((err) => {
  console.error('Failed to start server:', err);
  process.exit(1);
});