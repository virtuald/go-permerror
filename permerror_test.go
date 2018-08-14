package permerror_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	. "github.com/virtuald/go-permerror"
)

func AssertPermanence(t *testing.T, err error, typ TemporaryType) {
	assert.Equal(t, typ, IsTemporary(err))
}

func AssertErrorStuff(t *testing.T, cause error, err error, typ TemporaryType, msg string) {
	AssertPermanence(t, err, typ)

	assert.Equal(t, cause, errors.Cause(err))
	assert.Equal(t, msg, err.Error())
}

type isPermanent struct{}

func (*isPermanent) Error() string   { return "perm" }
func (*isPermanent) Temporary() bool { return false }

type isTemporary struct{}

func (*isTemporary) Error() string   { return "temp" }
func (*isTemporary) Temporary() bool { return true }

func TestIsTemporary(t *testing.T) {
	AssertPermanence(t, fmt.Errorf("nope"), Unknown)
	AssertPermanence(t, &isPermanent{}, Permanent)
	AssertPermanence(t, &isTemporary{}, Temporary)
}

func TestPermErrors(t *testing.T) {

	te := &isTemporary{}
	pe := &isPermanent{}
	un := fmt.Errorf("unknown")

	AssertErrorStuff(t, te, MakePermanent(te), Permanent, "temp")
	AssertErrorStuff(t, pe, MakePermanent(pe), Permanent, "perm")

	ne := New("yup")

	AssertErrorStuff(t, ne, ne, Permanent, "yup")

	AssertErrorStuff(t, un, WithMessage(un, "message"), Permanent, "message: unknown")
	AssertErrorStuff(t, te, WithMessage(te, "message"), Temporary, "message: temp")
	AssertErrorStuff(t, pe, WithMessage(pe, "message"), Permanent, "message: perm")

	AssertErrorStuff(t, un, Wrap(un), Permanent, "unknown")
	AssertErrorStuff(t, te, Wrap(te), Temporary, "temp")
	AssertErrorStuff(t, pe, Wrap(pe), Permanent, "perm")
}
