# Istio-authn-poc

Retrieve jwks:
```
curl -v http://localhost:8080/authn/v1/oauth2/jwks
```

Login:
```
curl -X POST http://localhost:8080/authn/v1/oauth2/token\?grant_type\=password\&scope\=everything\&username\=T00001\&password\=NtMiHEypfVJpIkPcDKYbAFAdfvTmaQQYiOpHelRJwAAmmqIDia\&client_id\=fake\&client_secret\=fakeSecret
```