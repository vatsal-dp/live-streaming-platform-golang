//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package message

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserControlEventDecodeCommon(t *testing.T) {
	for _, tc := range uceTestCases {
		tc := tc // capture

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewReader(tc.Binary)
			dec := NewUserControlEventDecoder(buf)

			var msg UserCtrlEvent
			err := dec.Decode(&msg)
			require.Nil(t, err)
			require.Equal(t, tc.Value, msg)
		})
	}
}
