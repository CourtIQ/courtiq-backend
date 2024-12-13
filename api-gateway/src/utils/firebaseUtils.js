// const admin = require('firebase-admin');
// const { FIREBASE_SERVICE_ACCOUNT } = require('../config');
// const { configDotenv } = require('dotenv');

// // Might be 
// if (!admin.apps.length) {
//   admin.initializeApp({
//     credential: admin.credential.cert(FIREBASE_SERVICE_ACCOUNT),
//   });
// }

// function verifyToken(idToken) {
//   // idToken comes from the client app
//   console.log('Verifying token:', idToken);
//   return getAuth()
//   .verifyIdToken(idToken)
//   .then((decodedToken) => {
//     const uid = decodedToken.uid;
//     console.log('Decoded token:', decodedToken);
//     return decodedToken;
//   })
//   .catch((error) => {
//     console.log('Error verifying token:', error);
//   });
//   // return admin.auth().verifyIdToken(idToken);
// }

// module.exports = { verifyToken };
