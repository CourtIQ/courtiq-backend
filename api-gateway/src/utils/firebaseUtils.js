const admin = require('firebase-admin');
const { FIREBASE_SERVICE_ACCOUNT } = require('../config');

if (!admin.apps.length) {
  admin.initializeApp({
    credential: admin.credential.cert(FIREBASE_SERVICE_ACCOUNT),
  });
}

async function verifyToken(idToken) {
  return admin.auth().verifyIdToken(idToken);
}

module.exports = { verifyToken };
