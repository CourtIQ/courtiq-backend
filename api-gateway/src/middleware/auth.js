const admin = require('firebase-admin');

if (!admin.apps.length) {
  // Parse the service account JSON from the environment variable
  const serviceAccount = JSON.parse(process.env.FIREBASE_SERVICE_ACCOUNT);

  admin.initializeApp({
    credential: admin.credential.cert(serviceAccount),
  });
}

async function authenticateToken(token) {
  if (!token) {
    throw new Error('Unauthorized: No token provided');
  }

  try {
    const decodedToken = await admin.auth().verifyIdToken(token);
    return decodedToken;
  } catch (error) {
    console.error('Authentication error:', error);
    throw new Error('Unauthorized: Invalid token');
  }
}

module.exports = authenticateToken;
