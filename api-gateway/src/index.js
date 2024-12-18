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

async function getIdTokenFromUid(uid) {
  const apiKey = "AIzaSyCyY5VHDOCqJKhek8o-q-s6LvFJ6kMNueQ";
  // Generate a custom token for the given UID
  const customToken = await admin.auth().createCustomToken(uid);

  // Identity Toolkit endpoint for signing in with a custom token
  const endpoint = `https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=${apiKey}`;

  const response = await fetch(endpoint, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      token: customToken,
      returnSecureToken: true,
      refreshToken: false
    })
  });

  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to exchange custom token: ${errorText}`);
  }

  const data = await response.json();
  console.log('Custom token exchange response:', data.idToken);
  // data should contain an idToken field
  return data.idToken;
}


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