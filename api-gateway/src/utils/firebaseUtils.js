const { admin } = require('../middleware/auth'); // Import the initialized admin instance
const axios = require('axios');
const config = require('../config');

async function exchangeCustomTokenForIdToken(customToken) {
  try {
    const response = await axios.post(
      `https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=${config.FIREBASE_CONFIG.apiKey}`,
      { token: customToken, returnSecureToken: true }
    );
    return response.data.idToken;
  } catch (error) {
    console.error('Error exchanging custom token for ID token:', error.message);
    throw error;
  }
}

async function generateIdToken(uid) {
  try {
    const customToken = await admin.auth().createCustomToken(uid); // Use the same initialized admin
    return exchangeCustomTokenForIdToken(customToken);
  } catch (error) {
    console.error('Error generating ID token:', error.message);
    throw error;
  }
}

module.exports = { exchangeCustomTokenForIdToken, generateIdToken };
