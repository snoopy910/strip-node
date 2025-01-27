# Start Sequencer

First run the sequencer DB:

```
docker run --name sequencer-postgres -p 5434:5432 -e POSTGRES_PASSWORD=password -d postgres
```

Then start the sequencer using:

```
go run main.go --httpPort="8082" -postgresHost="localhost:5434" --isSequencer=true
```

# Status

Intent Status changes in the following way:

```
Processing -> Failed -> Completed
```

Operations Status changes in the following way:

```
Pending -> Waiting -> Failed or Completed
```

Whenever an intent is created a goroutine is created to handle the intent. It will process the intent operations one by one. In case process is exited while an intent is ongoing then after process is started all the process intents are handled. 

The goroutine that handles the intents, will loop through the operations and handle them one by one.

# Sequener APIs

### Create Wallet

```
curl http://localhost:8082/createWallet?identity=0x40310390beF6D518f62Cc41a60a4E9b4a8b98730&identityCurve=ecdsa
```

### Get Wallet

```
curl http://localhost:8082/getWallet?identity=0x40310390beF6D518f62Cc41a60a4E9b4a8b98730&identityCurve=ecdsa
```

### Create Intent

```
curl --location 'localhost:8082/createIntent' \
--header 'Content-Type: application/json' \
--data '{
    "signature": "0x01",
    "identity": "0x40310390beF6D518f62Cc41a60a4E9b4a8b98730",
    "identityCurve": "ecdsa",
    "operations": [{
        "serializedTxn": "eb808477359400825208941c3bb7529e6a484db8f6d7f30e8e840c68dcd13788016345785d8a000080808080",
        "dataToSign": "cdf74de9e6b07952da2e4c6f38e14c1536abefeb9e1f37240e9392e195785c20",
        "chainId": "1337",
        "keyCurve": "ecdsa"
    },{
        "serializedTxn": "187PYuMkJ1sTeui7KoDted7fTc7BimkvmbK5VR4f6vwCBdotDjKZbMtmNn976iWqPoqFesmMEKDdj3SqmgWvQQhD4wsuov477Wd6r2yjinZ3BiBKupWYTvi2U9iqq6eRTUnKKdJV7rkywHeHimYR1jtAqzjsrjM9gt6aWppRjjULx94PCYqZbkEcbp3xeyHT7aBmPMAU2xRg3",
        "dataToSign": "187PYuMkJ1sTeui7KoDted7fTc7BimkvmbK5VR4f6vwCBdotDjKZbMtmNn976iWqPoqFesmMEKDdj3SqmgWvQQhD4wsuov477Wd6r2yjinZ3BiBKupWYTvi2U9iqq6eRTUnKKdJV7rkywHeHimYR1jtAqzjsrjM9gt6aWppRjjULx94PCYqZbkEcbp3xeyHT7aBmPMAU2xRg3",
        "chainId": "901",
        "keyCurve": "eddsa"
    }]
}'
```

### Get Intent by Id

```
curl --location 'localhost:8082/getIntent?id=1'
```

### Get Intents by Status

```
curl --location 'localhost:8082/getIntents?status=processing'
```

# Google oauth2 login


### Setup



Steps to setup Google oauth2 login require the following:

* Creating the credentials - Client ID and Client Secret - described in:

 https://developers.google.com/workspace/guides/create-credentials

* Creating the oauth consent screen - described in:

 https://developers.google.com/workspace/guides/configure-oauth-consent

Example for an oauth instance for StripChain:

Credentials: https://console.cloud.google.com/apis/credentials?inv=1&invt=AbnkgA&project=stripchain

OAuth consent screen: https://console.cloud.google.com/apis/credentials/consent?inv=1&invt=AbnkgA&project=stripchain


### Endpoints

`/oauth/login`: Login with Google

`/oauth/callback`: Callback from Google - checks if the user is authenticated and have a valid access token otherwise it will authenticate the user. A wallet will be created for the user if it does not exist derived from the Google ID. This wallet will be used as the identity for the user when they login with the Google oauth.

`/oauth/sign`: signs the message with the google wallet derived identity.

`/oauth/accessToken`: Generates a new access token when the current refresh token is still valid. It will return the new access token and the new refresh token

`/oauth/logout`: Logs out the user

A middleware validator is used at the router level to check if the user is authenticated and have a valid access token. It is activated if oauth authentication is enabled by setting the environment variable `ENABLE_OAUTH=true` or the `--enableOauth` flag.

Endpoints `/oauth/*` are exempt from the middleware validator.

The middleware validator checks if the URL path contains the `auth` query parameter with the value `oauth`. If it does, the middleware will check if the user is authenticated and have a valid access token. If the `auth` query parameter with the value `oauth` is not present, the middleware will pass the request to the next handler without checking if the user is authenticated and have a valid access token (to be validated with Nikolay).

Example:

oauth enabled and auth=`oauth`: Middleware will check if the user is authenticated and have a valid access token

http://localhost/createWallet?identity=0x2c8251052663244f37BAc7Bde1C6Cb02bBffff93&identityCurve=ecdsa&auth=oauth

oauth enabled or not enabled and auth is not `oauth` will keep the behavior as before without access token verification.

http://localhost/createWallet?identity=0x2c8251052663244f37BAc7Bde1C6Cb02bBffff93&identityCurve=ecdsa 


### Environment Variables

* `CLIENT_ID` - The client ID of the Google OAuth 2.0 application.

* `CLIENT_SECRET` - The client secret of the Google OAuth 2.0 application.

* `REDIRECT_URL` - The redirect URI for the Google OAuth 2.0 application.

* `JWT_SECRET` - The secret key used to sign and verify JWT tokens.

* `SESSION_SECRET` - The secret key used to encrypt and decrypt session data.

* `SALT` - The salt used to generate deterministic private keys for identity derivation from Google ID.

* `ENABLE_OAUTH` - If set to true, the OAuth endpoints will be enabled.

