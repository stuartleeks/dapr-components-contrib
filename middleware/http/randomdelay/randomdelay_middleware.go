/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package randomdelay

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"time"

	contribMetadata "github.com/dapr/components-contrib/metadata"
	"github.com/dapr/components-contrib/middleware"

	"github.com/dapr/kit/logger"
	kitmd "github.com/dapr/kit/metadata"
)

type randomDelayMiddlewareMetadata struct {
	MaxDelayMs int `json:"maxDelayMs"`
}

const (
	maxDelayMsKey = "maxDelayMs"

	// Defaults.
	defaultMaxDelayMs = 2000
)

type Middleware struct{}

func NewRandomDelayMiddleware(_ logger.Logger) middleware.Middleware {
	return &Middleware{}
}

func (m *Middleware) GetHandler(_ context.Context, metadata middleware.Metadata) (func(next http.Handler) http.Handler, error) {
	meta, err := m.getNativeMetadata(metadata)
	if err != nil {
		return nil, err
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Get random value between 0 and MaxDelayMs
			delayMs := rand.Intn(meta.MaxDelayMs)
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
			next.ServeHTTP(w, r)
		})
	}, nil
}

func (m *Middleware) getNativeMetadata(metadata middleware.Metadata) (*randomDelayMiddlewareMetadata, error) {
	middlewareMetadata := randomDelayMiddlewareMetadata{
		MaxDelayMs: defaultMaxDelayMs,
	}
	err := kitmd.DecodeMetadata(metadata.Properties, &middlewareMetadata)
	if err != nil {
		return nil, err
	}

	if middlewareMetadata.MaxDelayMs <= 0 {
		return nil, fmt.Errorf("metadata property %s must be a positive value", maxDelayMsKey)
	}

	return &middlewareMetadata, nil
}

func (m *Middleware) GetComponentMetadata() (metadataInfo contribMetadata.MetadataMap) {
	metadataStruct := randomDelayMiddlewareMetadata{}
	contribMetadata.GetMetadataInfoFromStructType(reflect.TypeOf(metadataStruct), &metadataInfo, contribMetadata.MiddlewareType)
	return metadataInfo
}
