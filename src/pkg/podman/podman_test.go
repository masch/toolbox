package podman_test

import (
	"errors"
	"github.com/containers/toolbox/pkg/podman"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestCheckVersion(t *testing.T) {
	type input struct {
		requireVersion string
	}

	type expect struct {
		result bool
	}

	tt := []struct {
		name   string
		input  input
		expect expect
	}{
		{
			name: "RequireVersion_GreaterThan_Supported_ShouldReturn_False",
			input: input{
				requireVersion: "10.1.1",
			},
			expect: expect{
				result: false,
			},
		},
		{
			name: "RequireVersion_LowerThan_Supported_ShouldReturn_True",
			input: input{
				requireVersion: "1.0.0",
			},
			expect: expect{
				result: true,
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expect.result,  podman.CheckVersion(tc.input.requireVersion))
		})
	}
}

func TestContainerExists(t *testing.T) {
	type input struct {
		containerName string
	}

	type expect struct {
		exists bool
		err error
	}

	tt := []struct {
		name   string
		input  input
		expect expect
	}{
		{
			name: "NonExisting_ShouldReturn_Error",
			input: input{
				containerName: "container-1",
			},
			expect: expect{
				exists: false,
				err:    errors.New("failed to find container container-1"),
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actualExists, actualErr := podman.ContainerExists(tc.input.containerName)
			require.EqualValues(t, tc.expect.err, actualErr)
			require.EqualValues(t, tc.expect.exists, actualExists)
		})
	}
}

func TestGetContainers(t *testing.T) {
	type input struct {
		args []string
	}

	type expect struct {
		err error
	}

	tt := []struct {
		name   string
		input  input
		expect expect
	}{
		{
			name: "ValidInputArgs_ShouldReturn_NonEmptyValues",
			input: input{
				args: nil,
			},
			expect: expect{
				err:    nil,
			},
		},
		{
			name: "InvalidInputArgs_ShouldReturn_NonEmptyValues",
			input: input{
				args: []string{"invalid"},
			},
			expect: expect{
				err:    errors.New("failed to invoke podman(1)"),
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actualContainers, actualErr := podman.GetContainers(tc.input.args...)
			require.EqualValues(t, tc.expect.err, actualErr)
			if tc.expect.err == nil {
				require.NotEmpty(t, actualContainers)
			} else {
				require.Empty(t, actualContainers)
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	type input struct {}

	type expect struct {
		version string
		err error
	}

	tt := []struct {
		name   string
		input  input
		expect expect
	}{
		{
			name: "ValidResponse_ShouldReturn_AValidVersionFormat",
			input: input{},
			expect: expect{
				version: "a",
				err:    nil,
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actualVersion, actualErr := podman.GetVersion()
			require.EqualValues(t, tc.expect.err, actualErr)
			if tc.expect.err == nil {
				// Since the actual version obtained is dynamic, we cannot ensure a fixed version,
				// so instead we check if the value has 3 values split with a dot format (x.x.x)
				require.EqualValues(t, 3, len(strings.Split(actualVersion, ".")))
			} else {
				require.NotEmpty(t,actualVersion)
			}
		})
	}
}