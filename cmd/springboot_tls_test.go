package cmd

import "testing"

// The verdict arrives on the same combined stream as compose's own container
// lifecycle chatter, so parsing has to find the sentinel line rather than read
// the output as a whole. Getting this wrong fails in the expensive direction:
// an unrecognised verdict is certUnknown, which lets issuance proceed and can
// spend one of the five duplicate certificates allowed per 168 hours.
func TestParseCertState(t *testing.T) {
	for _, tc := range []struct {
		name       string
		out        string
		want       certState
		wantDetail string
	}{
		{
			name: "live certificate surrounded by compose output",
			out: "[+] Creating 1/0\n" +
				" ✔ Container splm-mis-certbot  Created\n" +
				"SERVKUN-CERT live Nov 18 09:14:22 2026 GMT\n",
			want:       certLive,
			wantDetail: "Nov 18 09:14:22 2026 GMT",
		},
		{name: "no certificate", out: "SERVKUN-CERT none", want: certNone},
		{name: "placeholder from a failed bootstrap", out: "SERVKUN-CERT placeholder", want: certPlaceholder},
		{name: "expired", out: "SERVKUN-CERT expired", want: certExpired},
		{
			// A pull that fails, a volume that won't mount: no verdict is
			// printed at all. Must not be mistaken for "no certificate", which
			// would report a healthy stack as needing issuance.
			name: "no verdict printed",
			out:  "Error response from daemon: pull access denied",
			want: certUnknown,
		},
		{name: "empty output", out: "", want: certUnknown},
		{
			// docker compose run over a TTY-less stream still leaves \r behind
			// on some hosts; a trailing \r must not become part of the verdict.
			name:       "carriage returns are trimmed",
			out:        "SERVKUN-CERT live Nov 18 09:14:22 2026 GMT\r\n",
			want:       certLive,
			wantDetail: "Nov 18 09:14:22 2026 GMT",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, detail := parseCertState(tc.out)
			if got != tc.want {
				t.Errorf("parseCertState(%q) state = %v, want %v", tc.out, got, tc.want)
			}
			if detail != tc.wantDetail {
				t.Errorf("parseCertState(%q) detail = %q, want %q", tc.out, detail, tc.wantDetail)
			}
		})
	}
}

// The certLive branch is the only one that refuses to act, and its message is
// the operator's only clue about the rate limit they would otherwise spend. A
// verdict with no expiry date behind it would print that message with a blank
// where the date belongs.
func TestParseCertStateLiveWithoutExpiryIsStillLive(t *testing.T) {
	got, detail := parseCertState("SERVKUN-CERT live")
	if got != certLive {
		t.Errorf("state = %v, want certLive", got)
	}
	if detail != "" {
		t.Errorf("detail = %q, want empty", detail)
	}
}
