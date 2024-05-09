package helper

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"math/big"
	"net/http"
	"os"
)

var (
	AwsRegion            = os.Getenv("AWS_REGION")
	AwsCognitoUserPoolID = os.Getenv("AWS_COGNITO_USER_POOL_ID")
)

type JWKKey struct {
	Keys []struct {
		Alg string `json:"alg"`
		E   string `json:"e"`
		N   string `json:"n"`
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		Use string `json:"use"`
		X   string `json:"x"`
		Y   string `json:"y"`
		Crv string `json:"crv"`
	} `json:"keys"`
}

func findRSAPublicKey(jwkKey JWKKey, kid string) (*rsa.PublicKey, error) {
	for _, key := range jwkKey.Keys {
		if key.Kid == kid {
			if key.Kty == "RSA" {
				eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
				if err != nil {
					return nil, err
				}
				nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
				if err != nil {
					return nil, err
				}

				exponent := big.NewInt(0)
				exponent.SetBytes(eBytes)
				modulus := big.NewInt(0)
				modulus.SetBytes(nBytes)

				publicKey := &rsa.PublicKey{
					N: modulus,
					E: int(exponent.Int64()),
				}

				return publicKey, nil
			}
		}
	}
	return nil, fmt.Errorf("key not found")
}

func fetchJWKS() (*JWKKey, error) {
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", AwsRegion, AwsCognitoUserPoolID)
	response, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var jwks JWKKey
	if err := json.NewDecoder(response.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}

func VerifyToken(tokenString string, attribute string) (interface{}, error) {
	jwks, err := fetchJWKS()
	if err != nil {
		log.Printf("Error fetching JWKS: %v", err)
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("could not find kid in token header")
		}

		publicKey, err := findRSAPublicKey(*jwks, kid)
		if err != nil {
			return nil, err
		}

		return publicKey, nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//log.Printf("Token is valid, %s: %s", attribute, claims[attribute])
		return claims[attribute], nil
	} else {
		log.Printf("Token is invalid")
		return nil, err
	}

}
