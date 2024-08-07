package lib

import (
	"crypto/sha256"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMerkleTreeCreation(t *testing.T) {
	var data []Content
	elements := []string{
		"a",
		"b",
		"c",
		"d",
		"e",
		"f",
	}

	for _, e := range elements {
		data = append(data, StringContent{val: e})
	}

	tree, err := NewTree(data)
	if err != nil {
		t.Fatal(err)
	}

	verified, err := tree.VerifyContent(StringContent{val: "c"})
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, true, verified)

	verified, err = tree.VerifyContent(StringContent{val: "x"})
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, false, verified)
}

type StringContent struct {
	val string
}

func (s StringContent) CalculateHash() ([]byte, error) {
	bytes := []byte(s.val)
	hash := sha256.Sum256(bytes)
	return hash[:], nil
}

func (s StringContent) Equals(other Content) (bool, error) {
	otherAsStringContent, ok := other.(StringContent)
	if !ok {
		return false, nil
	}
	return s.val == otherAsStringContent.val, nil
}
