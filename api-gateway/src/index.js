const { ApolloGateway, IntrospectAndCompose } = require('@apollo/gateway');
const { ApolloServer } = require('@apollo/server');
const { expressMiddleware } = require('@apollo/server/express4');
const express = require('express');
const cors = require('cors');
const admin = require('firebase-admin');
const fs = require('fs');
const path = require('path');
const { services } = require('./config');

// Debug: Log service account file presence
const serviceAccountPath = path.join(__dirname, '..', 'firebase-service-account.json');
console.log('Checking service account file...');
console.log('Service account path:', serviceAccountPath);
console.log('File exists:', fs.existsSync(serviceAccountPath));

try {
  // Read the file content directly
  const serviceAccountContent = fs.readFileSync(serviceAccountPath, 'utf8');
  const serviceAccount = JSON.parse(serviceAccountContent);
  
  console.log('Service account loaded. Project ID:', serviceAccount.project_id);
  
  // Initialize Firebase Admin SDK
  admin.initializeApp({
    credential: admin.credential.cert(serviceAccount)
  });
  
  console.log('Firebase Admin SDK initialized successfully');
} catch (error) {
  console.error('Error initializing Firebase:', error);
  process.exit(1);
}

async function generateCustomToken(uid) {
  try {
    const customToken = await admin.auth().createCustomToken(uid);
    return customToken;
  } catch (error) {
    console.error('Error creating custom token:', error);
    return null;
  }
}

const gateway = new ApolloGateway({
  supergraphSdl: new IntrospectAndCompose({
    subgraphs: services,
    pollIntervalInMs: 1000,
  }),
});

async function startServer() {
  const app = express();

  const server = new ApolloServer({
    gateway,
    introspection: true,
  });

  await server.start();

  // Enhanced CORS configuration
  app.use(cors({
    origin: '*',
    methods: 'GET,HEAD,PUT,PATCH,POST,DELETE',
    preflightContinue: false,
    optionsSuccessStatus: 204,
    credentials: true
  }));
  
  app.use(express.json());

  // GraphQL Playground HTML
  app.get('/playground', (req, res) => {
    res.send(`
      <!DOCTYPE html>
      <html>
        <head>
          <title>GraphQL Playground</title>
          <meta charset="utf-8">
          <meta name="viewport" content="user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, minimal-ui">
          <link href="https://unpkg.com/graphql-playground-react/build/static/css/index.css" rel="stylesheet">
          <script src="https://unpkg.com/graphql-playground-react/build/static/js/middleware.js"></script>
        </head>
        <body>
          <div id="root">
            <style>
              body {
                background-color: rgb(23, 42, 58);
                font-family: Open Sans, sans-serif;
                height: 90vh;
              }
              #root {
                height: 100%;
                width: 100%;
                display: flex;
                align-items: center;
                justify-content: center;
              }
              .loading {
                font-size: 32px;
                font-weight: 200;
                color: rgba(255, 255, 255, .6);
                margin-left: 20px;
              }
              img {
                width: 78px;
                height: 78px;
              }
              .title {
                font-weight: 400;
              }
            </style>
            <img src='https://cdn.jsdelivr.net/npm/graphql-playground-react/build/logo.png' alt=''>
            <div class="loading"> Loading... </div>
          </div>
          <script>
            window.addEventListener('load', function (event) {
              const root = document.getElementById('root');
              root.classList.add('playgroundIn');
              const wsProto = location.protocol == 'https:' ? 'wss:' : 'ws:';
              GraphQLPlayground.init(root, {
                endpoint: '/graphql',
                subscriptionsEndpoint: wsProto + '//' + location.host + '/graphql',
                settings: {
                  'request.credentials': 'same-origin',
                }
              })
            })
          </script>
        </body>
      </html>
    `);
  });

  // Current token endpoint
  app.get('/current-token', async (req, res) => {
    const token = await generateCustomToken('1TM5LfDHmDbcIQS1X0mlngIZXWv1');
    console.log('Generated token:', token);
    res.json({ token });
  });

  // Token Generator UI
  app.get('/', (req, res) => {
    res.send(`
      <html>
        <head>
          <title>GraphQL API Gateway</title>
          <style>
            body { font-family: Arial, sans-serif; margin: 20px; background: #f0f0f0; }
            .container { max-width: 800px; margin: 0 auto; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
            .token { word-break: break-all; background: #f8f8f8; padding: 10px; margin: 10px 0; border-radius: 4px; border: 1px solid #ddd; }
            button { padding: 10px 20px; background: #4CAF50; color: white; border: none; border-radius: 4px; cursor: pointer; }
            button:hover { background: #45a049; }
            .copy-btn { margin-left: 10px; background: #2196F3; }
            .copy-btn:hover { background: #1976D2; }
            .links { margin-top: 20px; }
            .links a { display: inline-block; margin-right: 10px; color: #2196F3; text-decoration: none; }
            .links a:hover { text-decoration: underline; }
          </style>
        </head>
        <body>
          <div class="container">
            <h1>GraphQL API Gateway</h1>
            <div class="links">
              <a href="/playground">📚 Open GraphQL Playground</a>
              <a href="/graphql">🚀 GraphQL Endpoint</a>
            </div>
            <hr>
            <h2>Token Generator</h2>
            <button onclick="generateToken()">Generate New Token</button>
            <div id="tokenContainer"></div>
            
            <script>
              async function generateToken() {
                const tokenContainer = document.getElementById('tokenContainer');
                tokenContainer.innerHTML = 'Generating token...';
                
                try {
                  const response = await fetch('/current-token');
                  const data = await response.json();
                  
                  if (!data.token) {
                    throw new Error('No token received');
                  }
                  
                  const token = 'Bearer ' + data.token;
                  tokenContainer.innerHTML = \`
                    <h3>Generated Token:</h3>
                    <div class="token" id="tokenText">\${token}</div>
                    <button onclick="copyToken()" class="copy-btn">Copy Token</button>
                    <p>Copy this token and use it in the GraphQL Playground headers:</p>
                    <pre>
{
  "authorization": "\${token}"
}
                    </pre>
                  \`;
                } catch (error) {
                  console.error('Error:', error);
                  tokenContainer.innerHTML = \`<p style="color: red">Error: \${error.message}</p>\`;
                }
              }
              
              function copyToken() {
                const tokenText = document.querySelector('.token').textContent;
                navigator.clipboard.writeText(tokenText).then(() => {
                  alert('Token copied to clipboard!');
                });
              }
            </script>
          </div>
        </body>
      </html>
    `);
  });

  app.use(
    '/graphql',
    expressMiddleware(server, {
      context: async ({ req }) => {
        // Temporarily allow all requests without authentication
        return {};  // Empty context since authentication is bypassed for now
      },
    })
  );
  
  // GraphQL endpoint
// Replace only the GraphQL endpoint middleware section with this updated version:
// app.use(
//   '/graphql',
//   expressMiddleware(server, {
//     context: async ({ req }) => {
//       // Allow introspection queries without auth
//       const isIntrospectionQuery = req.body?.query?.includes('__schema') || 
//                                  req.body?.query?.includes('__type');
      
//       if (isIntrospectionQuery) {
//         return {};
//       }

//       const authHeader = req.headers.authorization;
//       if (!authHeader || !authHeader.startsWith('Bearer ')) {
//         throw new Error('No valid authorization header provided');
//       }

//       try {
//         const token = authHeader.split('Bearer ')[1];
        
//         // Get user from the token
//         // We'll use the test UID for development
//         const uid = 'UrO93Cetl8RMbHxUWGnE78V29rj2';
//         const userRecord = await admin.auth().getUser(uid);

//         return {
//           user: { uid },
//           userRecord,
//         };
//       } catch (error) {
//         console.error('Authentication error:', error);
//         throw new Error(`Authentication failed: ${error.message}`);
//       }
//     },
//   })
// );

  const PORT = process.env.PORT || 4000;
  app.listen(PORT, () => {
    console.log(`\n🚀 Gateway ready at http://localhost:${PORT}`);
    console.log(`📚 GraphQL Playground available at http://localhost:${PORT}/playground`);
    console.log(`⚡️ GraphQL endpoint available at http://localhost:${PORT}/graphql`);
  });
}

startServer().catch((err) => {
  console.error('Failed to start the server:', err);
  process.exit(1);
});