const { RemoteGraphQLDataSource } = require('@apollo/gateway');

class AuthenticatedDataSource extends RemoteGraphQLDataSource {
  willSendRequest({ request, context }) {
    if (context.user) {
      request.http.headers.set('x-user-id', context.user.uid);
    }
  }
}

module.exports = AuthenticatedDataSource;