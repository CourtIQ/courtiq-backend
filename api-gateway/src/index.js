require('dotenv').config();
const { ApolloServer } = require('@apollo/server');
const { startStandaloneServer } = require('@apollo/server/standalone');
const admin = require('firebase-admin');
const config = require('./config');
const gateway = require('./gateway');
const { refreshToken } = require('firebase-admin/app');

admin.initializeApp({
  credential: admin.credential.cert(config.FIREBASE_SERVICE_ACCOUNT),
});

async function startServer() {
  const server = new ApolloServer({
    gateway,
    introspection: config.GRAPHQL_PLAYGROUND,
  });

  const { url } = await startStandaloneServer(server, {
    listen: { port: config.PORT },
    context: async ({ req }) => {
      const authHeader = req.headers.authorization || '';
      const token = authHeader.replace('Bearer ', '');
      let user = null;
      
      if (token) {
        try {
          user = await admin.auth().verifyIdToken(token);
        } catch (error) {
          console.error('Error verifying token:', error);
        }
      }

      return { user };
    },
  });

  console.log(`🚀 Server ready at ${url}`);
}

startServer().catch((err) => {
  console.error('Failed to start server:', err);
  process.exit(1);
});