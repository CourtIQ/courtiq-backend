// middleware/auth.js
const admin = require('firebase-admin');
const config = require('../config');

// Initialize Firebase Admin
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

// Authentication middleware
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

module.exports = {
  admin,
  authenticateToken,
};