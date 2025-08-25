package interceptor

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func newTestUnaryHandler(t *testing.T, req any, err error) grpc.UnaryHandler {
	t.Helper()

	return func(_ context.Context, _ any) (any, error) {
		return req, err
	}
}

type mockValidator struct {
	err error
}

func (m *mockValidator) Validate() error {
	return m.err
}

type mockValidationFieldError struct {
	field, reason string
	cause         error
}

func (m *mockValidationFieldError) Error() string {
	return m.field + ": " + m.reason
}

func (m *mockValidationFieldError) Field() string {
	return m.field
}

func (m *mockValidationFieldError) Reason() string {
	return m.reason
}

func (m *mockValidationFieldError) Cause() error {
	return m.cause
}

func Test_Validate_NoValidator_ExpectOk(t *testing.T) {
	t.Parallel()

	got, err := Validation(context.Background(), nil, nil, newTestUnaryHandler(t, "test", nil))
	require.Equal(t, "test", got)
	require.NoError(t, err)
}

func Test_Validate_Validator_ExpectOk(t *testing.T) {
	t.Parallel()

	got, err := Validation(context.Background(),
		&mockValidator{errors.New("test mock validator err")},
		nil,
		newTestUnaryHandler(t, "test unary handler", nil),
	)
	require.Empty(t, got)
	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = Validation failed: test mock validator err")
}

func Test_Validate_ValidatorWithFieldError_ExpectOk(t *testing.T) {
	t.Parallel()

	got, err := Validation(context.Background(),
		&mockValidator{
			err: &mockValidationFieldError{
				field:  "someField1",
				reason: "someReason1",
			},
		},
		nil,
		newTestUnaryHandler(t, "test unary handler", nil),
	)
	require.Empty(t, got)
	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = Validation failed: someField1 someReason1")
}

func Test_Validate_ValidatorWithFieldErrorWithCause_ExpectOk(t *testing.T) {
	t.Parallel()

	got, err := Validation(context.Background(),
		&mockValidator{
			err: &mockValidationFieldError{
				field:  "someField1",
				reason: "someReason1",
				cause: &mockValidationFieldError{
					field:  "nestedField1",
					reason: "nestedReason1",
				},
			},
		},
		nil,
		newTestUnaryHandler(t, "test unary handler", nil),
	)
	require.Empty(t, got)
	require.EqualError(t, err, "rpc error: code = InvalidArgument desc = Validation failed: nestedField1 nestedReason1")
}
