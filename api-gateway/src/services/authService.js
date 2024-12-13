// const TokenVerifier = require('../utils/tokenVerifier');

// class AuthService {

//     constructor(tokenVerifier) {
//       this.tokenVerifier = new TokenVerifier();
//     }
  
//     async authenticate(headers) {
//       const authHeader = headers.authorization || '';
//       const token = authHeader.startsWith('Bearer ') ? authHeader.slice(7) : null;
  
//       if (!token) {
//         // return null;  // Allows requests without tokens to proceed as unauthenticated
//         throw new Error('No auth header provided'); // Strict error if no token is present
//       }
  
//       try {
//         console.log('Verifying token in authenticate:', token);
//         console.log('Token verifier:', this.tokenVerifier);
//         const userClaims = await this.tokenVerifier.verify(token);
//         console.log('User claims:', userClaims);
//         return userClaims;
//       } catch (error) {
//         // return null; // Silent fail, treats invalid tokens as unauthenticated
//         throw new Error('Invalid token in authenticate', error);
//       }
//   }
// }
  
//   module.exports = AuthService;
  