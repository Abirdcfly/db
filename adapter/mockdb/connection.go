// Copyright (c) 2012-today The upper.io/db authors. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package mockdb

import (
	"fmt"
	"net/url"
	"strings"
)

const connectionScheme = `mockdb`

type ConnectionURL struct {
	Database string
	Options  map[string]string
}

func (c ConnectionURL) String() (s string) {
	vv := url.Values{}

	if c.Options == nil {
		c.Options = map[string]string{}
	}

	if c.Database == "" {
		c.Database = "mockdb"
	}

	for k, v := range c.Options {
		vv.Set(k, v)
	}

	u := url.URL{
		Scheme:   connectionScheme,
		Path:     c.Database,
		RawQuery: vv.Encode(),
	}

	return u.String()
}

// ParseURL parses s into a ConnectionURL struct.
func ParseURL(s string) (conn ConnectionURL, err error) {
	var u *url.URL

	if !strings.HasPrefix(s, connectionScheme+"://") {
		return conn, fmt.Errorf("expecting %s:// connection scheme, got: %q", connectionScheme, s)
	}

	if u, err = url.Parse(s); err != nil {
		return conn, err
	}

	conn.Database = u.Host + u.Path
	conn.Options = map[string]string{}

	var vv url.Values

	if vv, err = url.ParseQuery(u.RawQuery); err != nil {
		return conn, err
	}

	for k := range vv {
		conn.Options[k] = vv.Get(k)
	}

	if _, ok := conn.Options["cache"]; !ok {
		conn.Options["cache"] = "shared"
	}

	return conn, err
}
