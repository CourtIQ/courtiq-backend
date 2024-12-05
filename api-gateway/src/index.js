// index.js
require('dotenv').config();
const { ApolloServer } = require('@apollo/server');
const { startStandaloneServer } = require('@apollo/server/standalone');
const { admin } = require('./middleware/auth');
const config = require('./config');
const gateway = require('./gateway');

async function startServer() {
  try {
    const server = new ApolloServer({
      gateway,
      introspection: true,
    });

    const { url } = await startStandaloneServer(server, {
      context: async ({ req }) => {
        const token = req.headers.authorization?.split(' ')[1] || null;
        return { token };
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