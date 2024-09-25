const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');
const { ApolloServer } = require('@apollo/server');
const { expressMiddleware } = require('@apollo/server/express4');
const express = require('express');
const cors = require('cors');
const admin = require('firebase-admin');
const { services } = require('./config');

// Parse the FIREBASE_CONFIG environment variable
const firebaseConfig = JSON.parse(process.env.FIREBASE_CONFIG);
var serviceAccount = require("../firebase-service-account.json");

// Initialize Firebase Admin SDK
admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
});

// Function to generate a custom Firebase JWT token for a specific user UID
async function generateTokenForUser(uid) {
  try {
    const customToken = await admin.auth().createCustomToken(uid);
    console.log(`Custom token for UID ${uid}:`, customToken);  // Moved the log here
    return customToken;
  } catch (error) {
    console.error('Error creating custom token:', error);
  }
}

// Uncomment the following line to generate a JWT token for the given UID
generateTokenForUser('1TM5LfDHmDbcIQS1X0mlngIZXWv1');  // Call the function here

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: services,
  }),
});

async function startServer() {
  const app = express();
  
  const server = new ApolloServer({
    gateway,
  });

  await server.start();
  const token = await generateTokenForUser('1TM5LfDHmDbcIQS1X0mlngIZXWv1');
  console.log(`Generated Token: ${token}`);
  app.use(cors());
  app.use(express.json());

  app.get('/', (req, res) => {
    res.send(`
      <html>
        <body>
          <h1>Welcome to the API Gateway</h1>
          <p>The GraphQL endpoint is available at <a href="/graphql">/graphql</a></p>
        </body>
      </html>
    `);
  });

  app.get('/health', (req, res) => {
    res.status(200).send('OK');
  });

  app.use(
    '/graphql',
    expressMiddleware(server, {
      context: async ({ req }) => {
        const token = req.headers.authorization?.split('Bearer ')[1];

        if (!token) {
          throw new Error('No token provided');
        }

        try {
          // Verify the Firebase JWT token
          const decodedToken = await admin.auth().verifyIdToken(token);
          
          // Check if the user exists in your Firebase project
          const userRecord = await admin.auth().getUser(decodedToken.uid);
          
          // Add the decoded token and user record to the context
          return { user: decodedToken, userRecord };
        } catch (error) {
          console.error('Authentication error:', error);
          if (error.code === 'auth/user-not-found') {
            throw new Error('User not found in this project');
          }
          throw new Error('Invalid token');
        }
      },
    })
  );

  const PORT = process.env.PORT || 4000;
  app.listen(PORT, () => {
    console.log(`🚀 Gateway ready at http://localhost:${PORT}`);
    console.log(`🚀 GraphQL endpoint available at http://localhost:${PORT}/graphql`);
  });
}

startServer().catch((err) => console.error('Failed to start the server:', err));
