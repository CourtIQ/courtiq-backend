const admin = require('firebase-admin');
const config = require('../config');

// Ensure Firebase Admin is initialized only once
if (!admin.apps.length) {
  try {
    admin.initializeApp({
      credential: admin.credential.cert(config.FIREBASE_SERVICE_ACCOUNT),
    });
    console.log('Firebase Admin initialized successfully.');
  } catch (error) {
    console.error('Failed to initialize Firebase Admin:', error.message);
    throw new Error('Firebase initialization failed');
  }
}

// Function to authenticate the ID token
async function authenticateToken(token) {
  if (!token) {
    throw new Error('Unauthorized: No token provided');
  }

  try {
    const decodedToken = await admin.auth().verifyIdToken(token);
    return decodedToken;
  } catch (error) {
    console.error('Authentication error:', error.message);
    throw new Error('Unauthorized: Invalid token');
  }
}

// GraphQL middleware to inject user into context
const authMiddleware = async (req, res, next) => {
  const authHeader = req.headers.authorization;

  if (!authHeader) {
    req.user = null; // No auth header provided
    return next();
  }

  const token = authHeader.split(' ')[1]; // Assumes "Bearer <token>"

  try {
    const decodedToken = await authenticateToken(token);
    req.user = {
      uid: decodedToken.uid,
      email: decodedToken.email,
      roles: decodedToken.roles || [], // If you have custom claims
    };
  } catch (error) {
    console.error('Error in authentication middleware:', error.message);
    req.user = null; // Optional: continue without user info
  }

  next();
};

module.exports = { authMiddleware };
