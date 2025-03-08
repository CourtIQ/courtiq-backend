const admin = require('firebase-admin');
const config = require('../config');
const logger = require('../logging/logger');

// Initialize Firebase Admin once at module level
let initialized = false;
const tokenCache = new Map(); // Simple token cache to reduce Firebase API calls
const TOKEN_CACHE_TTL = 5 * 60 * 1000; // 5 minutes in milliseconds

/**
 * Initializes the Firebase Admin SDK
 * @throws {Error} If initialization fails
 */
function initializeFirebase() {
  if (!initialized) {
    try {
      // Check if we have the required configuration
      if (!config.FIREBASE_SERVICE_ACCOUNT ||
          Object.keys(config.FIREBASE_SERVICE_ACCOUNT).length === 0) {
        logger.warn('Firebase service account not provided. Authentication will not work.');
        return;
      }
      
      admin.initializeApp({
        credential: admin.credential.cert(config.FIREBASE_SERVICE_ACCOUNT),
      });
      
      initialized = true;
      logger.info('Firebase Admin SDK initialized successfully');
    } catch (error) {
      logger.error('Failed to initialize Firebase Admin SDK', {
        error: error.message,
        stack: error.stack
      });
      throw error;
    }
  }
}

/**
 * Verifies a Firebase ID token with caching
 * @param {string} token - The Firebase ID token to verify
 * @returns {Promise<object|null>} The decoded token or null if invalid
 */
async function verifyIdToken(token) {
  if (!token) return null;
  
  try {
    // Make sure Firebase is initialized
    if (!initialized) {
      initializeFirebase();
      
      // If still not initialized after attempt, return null
      if (!initialized) {
        logger.warn('Firebase not initialized, cannot verify token');
        return null;
      }
    }
    
    // Check cache first
    const cachedData = tokenCache.get(token);
    if (cachedData && cachedData.expiresAt > Date.now()) {
      logger.debug('Using cached token validation');
      return cachedData.user;
    }
    
    // Verify token with Firebase
    const decodedToken = await admin.auth().verifyIdToken(token);
    
    // Cache the result with expiration
    const expiresAt = Date.now() + TOKEN_CACHE_TTL;
    tokenCache.set(token, { user: decodedToken, expiresAt });
    
    // Clean up expired cache entries periodically
    if (tokenCache.size > 100) { // Simple size-based cleanup trigger
      cleanupTokenCache();
    }
    
    return decodedToken;
  } catch (error) {
    logger.warn('Error verifying Firebase token', {
      error: error.message,
      errorCode: error.code,
      tokenPresent: !!token
    });
    return null;
  }
}

/**
 * Clean up expired tokens from the cache
 */
function cleanupTokenCache() {
  const now = Date.now();
  let expiredCount = 0;
  
  for (const [key, value] of tokenCache.entries()) {
    if (value.expiresAt <= now) {
      tokenCache.delete(key);
      expiredCount++;
    }
  }
  
  if (expiredCount > 0) {
    logger.debug(`Cleaned up ${expiredCount} expired tokens from cache`);
  }
}

module.exports = {
  initializeFirebase,
  verifyIdToken
};