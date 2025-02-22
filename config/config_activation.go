package config

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"r3/log"
	"r3/types"
	"slices"
)

// license revocations
var revocations = []string{"LI00334231"}

// License activation
// public key of custom key, used to verify license signature, created: 2025-02-19
// private key is stored in a secure location
// This is experimental, do not use in production
// Please use the upstream public key if use in production
// Download R3 at https://rei3.de/en/download
// You should contact https://rei3.de/en/contact for professional support
// If you have questions about this fork, please contact me: allan.registos@gmail.com

var publicKey = `-----BEGIN RSA PUBLIC KEY-----
MIICCgKCAgEAr6OcDoMUeg9HC/YgCl4lp9dp95YxWpRbcyVvCD/xHX3ThP0AbtFV
UNo3M/XOpn8xHUVHQSdzFArxOc3Flk+szirsDKeS3j6PDRWTWVfJtWaP5xX5iWs1
aUnLzX+SwFSfZLRglA92IyiShg3cM7cb2ZpewDnHfVfsnl9zEiOJc5vUDSQo2dOV
aygqPfeSP+IyaHo8AnXF2QF3BkcO9P4RaKN+nzkHMRQ91UrfuAK6w7X8y8kk3ytI
OI52Sm0LaM0Ce97ZzQg+snyQgr9Nmn7ODyi8lkNLwiq1qb0aY81vZAgULJosbwrf
+fFAARodWYQJDxxH0e3qF2MAJFZSwET6NAzUGK3KCxzoB9hNh2ki7NdhRTdziKcf
PcJs2bMmLkrQX17Elq263O7Pr33lUcS93sIzUtmqiimcDG+Y1Smb2SPk3XP3/Sdr
uHhGd7ZCf7MDe2ZCikQym0B4oSK5ChdIt7P4GD/Et0zY0PMa8311gTnI2JSgH48Z
k6f/HwAW+uIk07yKKwyk31/0a80wieS0rLwW2kcV+GrgrRGCvRRDjNcOjtRqXgs/
q1uhVClBgT6LsdzMjwd3SleM8WxwWVBhhoJmsOH0xwnTPQ2BtfA0qy+X/hU0VbIV
xMnFqSZdXyLz69vW24mTT9QCA61TwXWsP9l3UnGzFGCoPhFveJYXU90CAwEAAQ==
-----END RSA PUBLIC KEY-----`

func ActivateLicense() {
	if GetString("licenseFile") == "" {
		log.Info("server", "skipping activation check, no license installed")

		// set empty in case license was removed
		SetLicense(types.License{})
		return
	}

	var licFile types.LicenseFile

	if err := json.Unmarshal([]byte(GetString("licenseFile")), &licFile); err != nil {
		log.Error("server", "could not unmarshal license from config", err)
		return
	}

	licenseJson, err := json.Marshal(licFile.License)
	if err != nil {
		log.Error("server", "could not marshal license data", err)
		return
	}
	hashed := sha256.Sum256(licenseJson)

	// get license signature
	signature, err := base64.URLEncoding.DecodeString(licFile.Signature)
	if err != nil {
		log.Error("server", "could not decode license signature", err)
		return
	}

	// verify signature
	data, _ := pem.Decode([]byte(publicKey))
	if data == nil {
		log.Error("server", "could not decode public key", errors.New(""))
		return
	}
	key, err := x509.ParsePKCS1PublicKey(data.Bytes)
	if err != nil {
		log.Error("server", "could not parse public key", errors.New(""))
		return
	}

	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signature); err != nil {
		log.Error("server", "failed to verify license", err)
		return
	}

	// check if license has been revoked
	if slices.Contains(revocations, licFile.License.LicenseId) {
		log.Error("server", "failed to enable license", fmt.Errorf("license ID '%s' has been revoked", licFile.License.LicenseId))
		return
	}

	// set license
	log.Info("server", "setting license")
	SetLicense(licFile.License)
}
