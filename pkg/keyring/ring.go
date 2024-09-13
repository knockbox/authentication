package keyring

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/knockbox/authentication/pkg/utils"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"sync"
	"time"
)

// KeySet holds all the keys within the set. Provides functionality for rotating keys and more.
type KeySet struct {
	keys               jwk.Set
	curveTypes         []CurveType
	keyLifespanSeconds int
	jwtLifespanSeconds int
	activeKeys         map[string]*time.Timer
	expiringKeys       map[string]*time.Timer

	mu sync.RWMutex
	l  hclog.Logger
}

// NewSet creates a new KeySet and generates n ecdsa P521 keys.
func NewSet(keyLifespanSeconds, jwtLifespanSeconds int, l hclog.Logger) (*KeySet, error) {
	if keyLifespanSeconds < jwtLifespanSeconds {
		return nil, errors.New("invalid jwt key lifespan, keys will expire for valid jwt(s)")
	}

	return &KeySet{
		keys:               jwk.NewSet(),
		curveTypes:         nil,
		keyLifespanSeconds: keyLifespanSeconds,
		jwtLifespanSeconds: jwtLifespanSeconds,
		expiringKeys:       make(map[string]*time.Timer),
		activeKeys:         make(map[string]*time.Timer),
		mu:                 sync.RWMutex{},
		l:                  l,
	}, nil
}

// SetCurveTypes assigns the given elliptic.Curve(s) to the KeySet. These will be used when generating Keys.
func (k *KeySet) SetCurveTypes(types ...CurveType) {
	k.curveTypes = types
}

// Generate creates n keys and adds them to the set.
func (k *KeySet) Generate(n int) error {
	if k.curveTypes == nil || len(k.curveTypes) == 0 {
		return errors.New("no curve type(s) assigned to this set")
	}

	for i := 0; i < n; i++ {
		idx, err := utils.CryptoRandom(len(k.curveTypes))
		if err != nil {
			return err
		}

		// Generate the key
		curve := k.curveTypes[idx]
		rawKey, err := ecdsa.GenerateKey(curve.GetEllipticCurve(), rand.Reader)
		if err != nil {
			return errors.New("failed to generate ecdsa private key")
		}

		// Convert raw key to jwk
		key, err := jwk.FromRaw(rawKey)
		if err != nil {
			return errors.New("failed to create jwk from raw ecdsa private key")
		}

		// Assign the algorithm
		_ = key.Set(jwk.AlgorithmKey, curve.GetAlgorithm())

		// Set a random kid
		if err := key.Set(jwk.KeyIDKey, uuid.NewString()); err != nil {
			return errors.New("failed to set jwk kid")
		}

		// Set the use to signature
		if err := key.Set(jwk.KeyUsageKey, jwk.ForSignature); err != nil {
			return errors.New("failed to set jwk use")
		}

		// Add the key to the set
		if err := k.keys.AddKey(key); err != nil {
			return errors.New("failed to add jwk to set")
		}
		k.l.Debug("generated key", "kid", key.KeyID())

		signingLifespan := time.Now().Add(k.getSigningDuration()).Sub(time.Now())
		totalLifespan := time.Now().Add(k.getLifespanDuration()).Sub(time.Now())

		// Func for moving keys out of the signing pool
		k.activeKeys[key.KeyID()] = time.AfterFunc(signingLifespan, func() {
			k.mu.Lock()
			defer k.mu.Unlock()

			delete(k.activeKeys, key.KeyID())
			k.l.Info("active key expired", "kid", key.KeyID())

			if err := k.Generate(1); err != nil {
				k.l.Error("failed to replace active key", "kid", key.KeyID())
			}
		})

		// Func for moving keys out of the set
		k.expiringKeys[key.KeyID()] = time.AfterFunc(totalLifespan, func() {
			k.mu.Lock()
			defer k.mu.Unlock()

			delete(k.expiringKeys, key.KeyID())
			_ = k.keys.RemoveKey(key)
		})
	}

	return nil
}

// GetRandomKey returns a random jwk.Key that is active in the current set.
func (k *KeySet) GetRandomKey() jwk.Key {
	k.mu.RLock()
	defer k.mu.RUnlock()

	if len(k.activeKeys) == 0 {
		return nil
	}

	var validKeys []jwk.Key
	for kid := range k.activeKeys {
		key, exists := k.keys.LookupKeyID(kid)
		if !exists {
			k.l.Warn("expected key was not found in set", "kid", kid)
			continue
		}

		validKeys = append(validKeys, key)
	}

	idx, _ := utils.CryptoRandom(len(k.activeKeys))
	return validKeys[idx]
}

// GetKeyById returns the jwk.Key if one exists for the given kid, otherwise the second return value is false if the key
// is not found.
func (k *KeySet) GetKeyById(kid string) (jwk.Key, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	return k.keys.LookupKeyID(kid)
}

// RevokeKeyById removes a key from the set.
func (k *KeySet) RevokeKeyById(kid string) {
	k.mu.Lock()
	defer k.mu.Unlock()

	// Lookup the key and remove it if it exists.
	key, exists := k.keys.LookupKeyID(kid)
	if !exists {
		return
	}
	_ = k.keys.RemoveKey(key)

	// Stop the active timer.
	if aliveTimer, ok := k.activeKeys[kid]; ok {
		if !aliveTimer.Stop() {
			<-aliveTimer.C
		}
	}
	delete(k.activeKeys, kid)

	// Stop the expiring timer.
	if expiringTimer, ok := k.expiringKeys[kid]; ok {
		if !expiringTimer.Stop() {
			<-expiringTimer.C
		}
	}
	delete(k.expiringKeys, kid)
}

// GetPrivateKeySet returns all the jwk.Key(s) stored in the set containing the private key information.
func (k *KeySet) GetPrivateKeySet() *KeySetResponse {
	k.mu.RLock()
	defer k.mu.RUnlock()

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	var keys []jwk.Key
	for it := k.keys.Keys(ctx); it.Next(ctx); {
		pair := it.Pair()
		if pair == nil {
			continue
		}

		key, ok := pair.Value.(jwk.Key)
		if !ok {
			k.l.Debug("failed to iterate key", "pair", pair)
		}
		keys = append(keys, key)
	}

	return &KeySetResponse{Keys: keys}
}

// GetPublicKeySet returns all the jwk.Key(s) stored in the set without any private key information.
func (k *KeySet) GetPublicKeySet() *KeySetResponse {
	k.mu.RLock()
	defer k.mu.RUnlock()

	privateKeys := k.GetPrivateKeySet()

	var publicKeys []jwk.Key
	for _, pk := range privateKeys.Keys {
		key, err := pk.PublicKey()
		if err != nil {
			k.l.Debug("failed to convert key to public key", "kid", pk.KeyID())
		}

		publicKeys = append(publicKeys, key)
	}

	return &KeySetResponse{Keys: publicKeys}
}

// getSigningDuration returns the duration in seconds the key is valid for signing new jwt(s).
func (k *KeySet) getSigningDuration() time.Duration {
	return time.Duration(k.keyLifespanSeconds-k.jwtLifespanSeconds) * time.Second
}

// getLifespanDuration returns the total duration in seconds before the key will be removed from the keys set.
func (k *KeySet) getLifespanDuration() time.Duration {
	return time.Duration(k.keyLifespanSeconds) * time.Second
}
