package pattern

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected Expr
		err      error
	}{
		{
			"<foo> bar f <f>",
			Expr{capture("foo"), literals(" bar f "), capture("f")},
			nil,
		},
		{
			"<foo",
			Expr{literals("<foo")},
			nil,
		},
		{
			"<foo ><bar>",
			Expr{literals("<foo >"), capture("bar")},
			nil,
		},
		{
			"<>",
			Expr{literals("<>")},
			nil,
		},
		{
			"<_>",
			Expr{capture("_")},
			nil,
		},
		{
			"<1_>",
			Expr{literals("<1_>")},
			nil,
		},
		{
			`<ip> - <user> [<_>] "<method> <path> <_>" <status> <size> <url> <user_agent>`,
			Expr{capture("ip"), literals(" - "), capture("user"), literals(" ["), capture("_"), literals(`] "`), capture("method"), literals(" "), capture("path"), literals(" "), capture('_'), literals(`" `), capture("status"), literals(" "), capture("size"), literals(" "), capture("url"), literals(" "), capture("user_agent")},
			nil,
		},
	} {
		tc := tc
		actual, err := parseExpr(tc.input)
		if tc.err != nil || err != nil {
			require.Equal(t, tc.err, err)
			return
		}
		require.Equal(t, tc.expected, actual)
	}
}
