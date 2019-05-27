//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package hot_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteController_Open(t *testing.T) {
	t.Parallel()

	t.Run("open for correct pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := hot.NewWriteController()
		err := m.Open(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("multiple open for same pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := hot.NewWriteController()
		err := m.Open(ctx, 1)
		require.NoError(t, err)

		err = m.Open(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("try to open previous pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := hot.NewWriteController()
		err := m.Open(ctx, 2)
		require.NoError(t, err)

		err = m.Open(ctx, 1)
		require.Error(t, err)
		assert.Equal(t, hot.ErrWriteClosed, err)
	})
}

func TestWriteController_CloseAndWait(t *testing.T) {
	t.Parallel()

	t.Run("close correct pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := hot.NewWriteController()
		_ = m.Open(ctx, 1)
		err := m.CloseAndWait(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("multiple close for same pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := hot.NewWriteController()
		_ = m.Open(ctx, 1)
		err := m.CloseAndWait(ctx, 1)
		require.NoError(t, err)

		err = m.CloseAndWait(ctx, 1)
		require.Error(t, err)
		assert.Equal(t, hot.ErrWriteClosed, err)
	})

	t.Run("try to close incorrect pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := hot.NewWriteController()
		err := m.Open(ctx, 2)
		require.NoError(t, err)

		err = m.CloseAndWait(ctx, 1)
		require.Error(t, err)
		assert.Equal(t, hot.ErrWriteClosed, err)

		err = m.CloseAndWait(ctx, 3)
		require.Error(t, err)
		assert.Equal(t, hot.ErrWriteClosed, err)
	})
}

func TestWriteController_Begin(t *testing.T) {
	t.Parallel()

	t.Run("begin for not-opened pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := hot.NewWriteController()
		_, err := m.Begin(ctx, 1)
		require.Error(t, err)
		assert.Equal(t, hot.ErrWriteClosed, err)
	})

	t.Run("begin for closed pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := hot.NewWriteController()
		err := m.Open(ctx, 1)
		require.NoError(t, err)
		err = m.CloseAndWait(ctx, 1)
		require.NoError(t, err)

		_, err = m.Begin(ctx, 1)
		require.Error(t, err)
		assert.Equal(t, hot.ErrWriteClosed, err)
	})

	t.Run("begin for correct pulse", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := hot.NewWriteController()
		err := m.Open(ctx, 1)
		require.NoError(t, err)

		for i := 0; i < 1000; i++ {
			done, _ := m.Begin(ctx, 1)
			go func() {
				time.Sleep((time.Duration)(rand.Int31n(1000)+10) * time.Millisecond)
				done()
			}()
		}
		err = m.CloseAndWait(ctx, 1)
		require.NoError(t, err)
	})

	t.Run("begin while waiting pulse closing", func(t *testing.T) {
		t.Parallel()
		ctx := inslogger.TestContext(t)

		m := hot.NewWriteController()
		err := m.Open(ctx, 1)
		require.NoError(t, err)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			for i := 0; i < 1000; i++ {
				done, _ := m.Begin(ctx, 1)
				go func() {
					time.Sleep((time.Duration)(rand.Int31n(1000)) * time.Millisecond)
					done()
				}()
			}
			_ = m.CloseAndWait(ctx, 1)
			wg.Done()
		}()

		time.Sleep(time.Millisecond * 100)

		_, err = m.Begin(ctx, 1)

		wg.Wait()

		require.Error(t, err)
		assert.Equal(t, hot.ErrWriteClosed, err)
	})
}