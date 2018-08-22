package keyvalue_test

import (
	"testing"

	"github.com/kowala-tech/kcoin/notifications/keyvalue"
	"github.com/stretchr/testify/require"
)

func wrappedMemoryKV() keyvalue.Value {
	kv := keyvalue.NewMemoryKeyValue()
	return keyvalue.WrapKeyValue(kv, "test")
}

func Test_MemoryValueStorage_StoreString(t *testing.T) {
	v := wrappedMemoryKV()
	str := "Hello world"

	err := v.PutString(str)
	require.NoError(t, err)

	received, err := v.GetString()
	require.NoError(t, err)
	require.Equal(t, received, str)
}

func Test_MemoryValueStorage_StoreInt(t *testing.T) {
	v := wrappedMemoryKV()
	var n int64 = 42

	err := v.PutInt64(n)
	require.NoError(t, err)

	received, err := v.GetInt64()
	require.NoError(t, err)
	require.Equal(t, received, n)
}
