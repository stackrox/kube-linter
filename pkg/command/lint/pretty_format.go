package lint

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/term"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/run"
)

// Pretty output: findings grouped by object -> category -> check, with a
// proportional overview bar. Port of the kube-lint-pretty Bash wrapper.
//
// Environment:
//   NO_COLOR=1     no ANSI colors
//   FORCE_COLOR=1  colors even without a TTY (CI logs)
//   NO_UNICODE=1   ASCII fallback
//   KLP_FIX=1      show remediation once per check group

type prettyCategory int

const (
	catSecurity prettyCategory = iota + 1
	catReliability
	catResources
	catHygiene
	catOther
)

var prettyCategoryNames = [...]string{"", "Security", "Reliability", "Resources", "Hygiene", "Other"}

var prettyCategoryByCheck = map[string]prettyCategory{
	// Security
	"access-to-create-pods": catSecurity, "access-to-secrets": catSecurity, "cluster-admin-role-binding": catSecurity,
	"default-service-account": catSecurity, "docker-sock": catSecurity, "drop-net-raw-capability": catSecurity,
	"env-var-secret": catSecurity, "exposed-services": catSecurity, "host-ipc": catSecurity, "host-network": catSecurity,
	"host-pid": catSecurity, "no-read-only-root-fs": catSecurity, "privilege-escalation-container": catSecurity,
	"privileged-container": catSecurity, "privileged-ports": catSecurity, "read-secret-from-env-var": catSecurity,
	"run-as-non-root": catSecurity, "sensitive-host-mounts": catSecurity, "ssh-port": catSecurity,
	"unsafe-proc-mount": catSecurity, "unsafe-sysctls": catSecurity, "wildcard-in-rules": catSecurity,
	"writable-host-mount": catSecurity, "scc-deny-privileged-container": catSecurity, "non-isolated-pod": catSecurity,
	// Reliability
	"minimum-three-replicas": catReliability, "hpa-minimum-three-replicas": catReliability, "no-anti-affinity": catReliability,
	"no-liveness-probe": catReliability, "no-readiness-probe": catReliability, "no-rolling-update-strategy": catReliability,
	"dangling-service": catReliability, "mismatching-selector": catReliability, "non-existent-service-account": catReliability,
	"liveness-port": catReliability, "readiness-port": catReliability, "startup-port": catReliability,
	"dangling-horizontalpodautoscaler": catReliability, "dangling-ingress": catReliability, "dangling-networkpolicy": catReliability,
	"dangling-networkpolicypeer-podselector": catReliability, "dangling-servicemonitor": catReliability, "no-node-affinity": catReliability,
	// Resources
	"unset-cpu-requirements": catResources, "unset-memory-requirements": catResources,
	// Hygiene
	"latest-tag": catHygiene, "use-namespace": catHygiene, "deprecated-service-account-field": catHygiene,
	"no-extensions-v1beta": catHygiene, "restart-policy": catHygiene, "invalid-target-ports": catHygiene,
	"job-ttl-seconds-after-finished": catHygiene, "dnsconfig-options": catHygiene, "duplicate-env-var": catHygiene,
	"env-value-from": catHygiene, "priority-class-name": catHygiene, "sorted-keys": catHygiene, "schema-validation": catHygiene,
}

func prettyCategoryFor(check string) prettyCategory {
	if c, ok := prettyCategoryByCheck[check]; ok {
		return c
	}
	switch {
	case strings.HasPrefix(check, "pdb-"):
		return catReliability
	case strings.HasPrefix(check, "required-"):
		return catHygiene
	}
	return catOther
}

var prettyTitles = map[string]string{
	"access-to-create-pods":            "Role may create pods",
	"access-to-secrets":                "Role may access secrets",
	"cluster-admin-role-binding":       "cluster-admin role binding",
	"default-service-account":          "Default service account in use",
	"dangling-service":                 "Service without matching pods",
	"deprecated-service-account-field": "Deprecated serviceAccount field",
	"docker-sock":                      "Docker socket mounted",
	"drop-net-raw-capability":          "NET_RAW capability configuration",
	"env-var-secret":                   "Secret-like value in env var",
	"exposed-services":                 "Service exposed externally",
	"host-ipc":                         "Host IPC namespace enabled",
	"host-network":                     "Host network namespace enabled",
	"host-pid":                         "Host process namespace enabled",
	"latest-tag":                       "Image uses :latest",
	"minimum-three-replicas":           "Fewer than three replicas",
	"mismatching-selector":             "Selector does not match labels",
	"no-anti-affinity":                 "No pod anti-affinity",
	"no-liveness-probe":                "Liveness probe missing",
	"no-read-only-root-fs":             "Writable root filesystem",
	"no-readiness-probe":               "Readiness probe missing",
	"no-rolling-update-strategy":       "No rolling update strategy",
	"non-existent-service-account":     "Service account does not exist",
	"privilege-escalation-container":   "Privilege escalation allowed",
	"privileged-container":             "Privileged container",
	"privileged-ports":                 "Privileged port in use",
	"read-secret-from-env-var":         "Secret read from env var",
	"run-as-non-root":                  "Container may run as root",
	"sensitive-host-mounts":            "Sensitive host path mounted",
	"ssh-port":                         "SSH port exposed",
	"unsafe-proc-mount":                "Unsafe /proc mount",
	"unsafe-sysctls":                   "Unsafe sysctls",
	"unset-cpu-requirements":           "CPU requirements missing",
	"unset-memory-requirements":        "Memory requirements missing",
	"use-namespace":                    "Default namespace in use",
	"wildcard-in-rules":                "Wildcard in RBAC rules",
	"writable-host-mount":              "Writable host mount",
}

func prettyTitleFor(check string) string {
	if t, ok := prettyTitles[check]; ok {
		return t
	}
	parts := strings.Split(check, "-")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, ucfirst(p))
		}
	}
	return strings.Join(out, " ")
}

func ucfirst(s string) string {
	r, n := utf8.DecodeRuneInString(s)
	if n == 0 {
		return s
	}
	return string(unicode.ToUpper(r)) + s[n:]
}

var (
	reContainer = regexp.MustCompile(`[Cc]ontainer "([^"]*)"`)
	reImage     = regexp.MustCompile(`image, "([^"]*)"`)
	reCPU       = regexp.MustCompile(`cpu (request|limit)`)
	reMemory    = regexp.MustCompile(`memory (request|limit)`)
)

func prettyMessage(check, msg string) string {
	cq := "Container"
	if m := reContainer.FindStringSubmatch(msg); m != nil {
		cq = fmt.Sprintf("Container %q", m[1])
	}
	out := ""
	switch check {
	case "drop-net-raw-capability":
		switch {
		case strings.Contains(msg, "has ADD capability"):
			out = cq + " adds the forbidden NET_RAW capability."
		case strings.Contains(msg, "does not drop"):
			out = cq + " does not drop NET_RAW."
		}
	case "host-network":
		out = "Workload shares the host network namespace."
	case "host-pid":
		out = "Workload shares the host process namespace."
	case "host-ipc":
		out = "Workload shares the host IPC namespace."
	case "latest-tag":
		if m := reImage.FindStringSubmatch(msg); m != nil {
			out = cq + " uses " + m[1] + "; use an immutable, versioned image tag."
		}
	case "no-read-only-root-fs":
		out = cq + " does not use a read-only root filesystem."
	case "privilege-escalation-container":
		out = cq + " permits privilege escalation."
	case "privileged-container":
		out = cq + " runs in privileged mode."
	case "run-as-non-root":
		out = cq + " is not configured with runAsNonRoot."
	case "unset-cpu-requirements":
		if m := reCPU.FindStringSubmatch(msg); m != nil {
			out = cq + " has no CPU " + m[1] + "."
		}
	case "unset-memory-requirements":
		if m := reMemory.FindStringSubmatch(msg); m != nil {
			out = cq + " has no memory " + m[1] + "."
		}
	}
	if out == "" {
		m := strings.TrimPrefix(msg, "object ")
		m = strings.TrimPrefix(m, "resource ")
		out = ucfirst(m)
	}
	return out
}

// prettyStyle holds terminal-dependent settings.
type prettyStyle struct {
	width                              int
	bold, dim, red, green, cyan, reset string
	catColor                           [6]string
	ok, fail, ell                      string
	glyph                              [6]string
	showFix                            bool
}

func newPrettyStyle(out io.Writer) prettyStyle {
	var s prettyStyle
	if w, ok := out.(nopWriteCloser); ok {
		out = w.Writer
	}

	isTTY, fd := false, 0
	if f, ok := out.(*os.File); ok {
		fd = int(f.Fd())
		isTTY = term.IsTerminal(fd)
	}

	useColor := false
	if os.Getenv("NO_COLOR") == "" {
		if os.Getenv("FORCE_COLOR") != "" {
			useColor = true
		} else if isTTY && os.Getenv("TERM") != "" && os.Getenv("TERM") != "dumb" {
			useColor = true
		}
	}

	useUnicode := false
	if os.Getenv("NO_UNICODE") == "" {
		loc := ""
		for _, k := range []string{"LC_ALL", "LC_CTYPE", "LANG"} {
			if v := os.Getenv(k); v != "" {
				loc = v
				break
			}
		}
		l := strings.ToLower(loc)
		useUnicode = strings.Contains(l, "utf-8") || strings.Contains(l, "utf8")
	}

	if useColor {
		s.bold, s.dim, s.red, s.green, s.cyan, s.reset = "\x1b[1m", "\x1b[2m", "\x1b[31m", "\x1b[32m", "\x1b[2;36m", "\x1b[0m"
		s.catColor = [6]string{"", "\x1b[31m", "\x1b[33m", "\x1b[36m", "\x1b[35m", "\x1b[90m"}
	}
	if useUnicode {
		s.ok, s.fail, s.ell = "✓", "✗", "…"
		s.glyph = [6]string{"", "█", "▓", "▒", "░", "·"}
	} else {
		s.ok, s.fail, s.ell = "[OK]", "[FAIL]", "..."
		s.glyph = [6]string{"", "#", "=", "+", "-", "."}
	}

	cols := 0
	if isTTY {
		if w, _, err := term.GetSize(fd); err == nil {
			cols = w
		}
	}
	if cols <= 0 {
		if v, err := strconv.Atoi(os.Getenv("COLUMNS")); err == nil && v > 0 {
			cols = v
		} else {
			cols = 100
		}
	}
	s.width = min(max(cols, 40), 84)

	s.showFix = os.Getenv("KLP_FIX") == "1"
	return s
}

func strLen(s string) int { return utf8.RuneCountInString(s) }

func truncRunes(s string, n int) string {
	if n <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// wrap breaks text into at most maxLines lines of the given width; truncation only at the end, marked with ellipsis.
func (s prettyStyle) wrap(text string, width, maxLines int) []string {
	ellLen := strLen(s.ell)
	var lines []string
	line := ""
	for _, word := range strings.Fields(text) {
		if strLen(word) > width {
			word = truncRunes(word, width-ellLen) + s.ell
		}
		switch {
		case line == "":
			line = word
		case strLen(line)+1+strLen(word) <= width:
			line += " " + word
		default:
			lines = append(lines, line)
			line = word
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		last := lines[maxLines-1]
		if strLen(last)+ellLen > width {
			last = truncRunes(last, width-ellLen)
		}
		if n := len(last); n > 0 && strings.ContainsRune(" .,;:", rune(last[n-1])) {
			last = last[:n-1]
		}
		lines[maxLines-1] = last + s.ell
	}
	return lines
}

// lr prints left and right text right-aligned on one line, or on two lines if it doesn't fit.
func (s prettyStyle) lr(w io.Writer, indent int, lp, lf, rp, rf string) {
	ind := strings.Repeat(" ", indent)
	pad := s.width - indent - strLen(lp) - strLen(rp)
	if pad < 2 {
		fmt.Fprintf(w, "%s%s\n%s%s\n", ind, lf, ind, rf)
		return
	}
	fmt.Fprintf(w, "%s%s%s%s\n", ind, lf, strings.Repeat(" ", pad), rf)
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

type prettyRow struct {
	key, check, src, ns, name, api, kind, msg, rem string
	cat                                            prettyCategory
}

func newPrettyRow(r diagnostic.WithContext) prettyRow {
	row := prettyRow{
		check: r.Check,
		src:   r.Object.Metadata.FilePath,
		msg:   r.Diagnostic.Message,
		rem:   r.Remediation,
		cat:   prettyCategoryFor(r.Check),
	}
	if r.Object.K8sObject != nil {
		info := r.Object.GetK8sObjectName()
		row.ns, row.name = info.Namespace, info.Name
		row.api = info.GroupVersionKind.GroupVersion().String()
		row.kind = info.GroupVersionKind.Kind
	}
	row.key = row.src + "|" + row.ns + "/" + row.name + "|" + row.kind
	return row
}

func (row prettyRow) context() string {
	s := row.kind
	if api := strings.TrimPrefix(row.api, "/"); api != "" {
		if s != "" {
			s += " "
		}
		s += api
	}
	if row.ns != "" {
		s += " in " + row.ns
	}
	src := row.src
	if wd, err := os.Getwd(); err == nil {
		if rel, err := filepath.Rel(wd, src); err == nil && !strings.HasPrefix(rel, "..") {
			src = rel
		}
	}
	if s != "" {
		return s + ", from " + src
	}
	return "from " + src
}

func formatLintPretty(out io.Writer, data interface{}) error {
	result, ok := data.(run.Result)
	if !ok {
		return fmt.Errorf("pretty format: unexpected data type %T", data)
	}
	s := newPrettyStyle(out)
	W := s.width

	total := len(result.Reports)
	if total == 0 {
		_, err := fmt.Fprintf(out, "%s%s%s KubeLinter: no findings.\n", s.green, s.ok, s.reset)
		return err
	}

	rows := make([]prettyRow, 0, total)
	var catCount [6]int
	for _, r := range result.Reports {
		row := newPrettyRow(r)
		catCount[row.cat]++
		rows = append(rows, row)
	}
	// Object > category > check, stable (keeps kube-linter's order within a check).
	sort.SliceStable(rows, func(i, j int) bool {
		a, b := rows[i], rows[j]
		if a.key != b.key {
			return a.key < b.key
		}
		if a.cat != b.cat {
			return a.cat < b.cat
		}
		return a.check < b.check
	})

	objCount := map[string]int{}
	catInObj := map[string]int{}
	chkInObj := map[string]int{}
	files := map[string]struct{}{}
	var objOrder []string
	for _, r := range rows {
		if _, seen := objCount[r.key]; !seen {
			objOrder = append(objOrder, r.key)
		}
		objCount[r.key]++
		catInObj[r.key+"\x1f"+strconv.Itoa(int(r.cat))]++
		chkInObj[r.key+"\x1f"+r.check]++
		files[r.src] = struct{}{}
	}
	nobj, nfiles := len(objOrder), len(files)

	// Header
	var lp, lf string
	if nobj == 1 {
		lp, lf = rows[0].name, s.bold+rows[0].name+s.reset
	} else {
		lp = fmt.Sprintf("%d %s in %d %s", nobj, plural(nobj, "object", "objects"), nfiles, plural(nfiles, "file", "files"))
		lf = lp
	}
	fw := plural(total, "finding", "findings")
	s.lr(out, 0, lp, lf, fmt.Sprintf("%d %s", total, fw), fmt.Sprintf("%s%s%d%s %s", s.bold, s.red, total, s.reset, fw))
	if nobj == 1 {
		fmt.Fprintf(out, "%s%s%s\n", s.dim, rows[0].context(), s.reset)
	}
	fmt.Fprintln(out)

	// Bar: every category at least one cell, remainder goes to the largest.
	var cells [6]int
	big, sum := 1, 0
	for i := 1; i <= 5; i++ {
		if catCount[i] > 0 {
			cells[i] = max(catCount[i]*W/total, 1)
		}
		sum += cells[i]
		if cells[i] > cells[big] {
			big = i
		}
	}
	cells[big] += W - sum

	var bar, legend strings.Builder
	for i := 1; i <= 5; i++ {
		if catCount[i] == 0 {
			continue
		}
		bar.WriteString(s.catColor[i] + strings.Repeat(s.glyph[i], cells[i]) + s.reset)
		if legend.Len() > 0 {
			legend.WriteString("   ")
		}
		fmt.Fprintf(&legend, "%s%s%s %s %s%d%s", s.catColor[i], s.glyph[i], s.reset, prettyCategoryNames[i], s.bold, catCount[i], s.reset)
	}
	fmt.Fprintf(out, "%s\n%s\n", bar.String(), legend.String())

	// Findings
	bi := 0
	if nobj > 1 {
		bi = 2
	}
	var prevKey, prevChk string
	var prevCat prettyCategory
	var groupN int
	shownFix := false
	for _, r := range rows {
		if r.key != prevKey {
			prevKey, prevCat, prevChk = r.key, 0, ""
			if nobj > 1 {
				n := objCount[r.key]
				fmt.Fprintln(out)
				pl := plural(n, "finding", "findings")
				s.lr(out, 0, r.name, s.bold+r.name+s.reset, fmt.Sprintf("%d %s", n, pl), fmt.Sprintf("%s%s%d%s %s", s.bold, s.red, n, s.reset, pl))
				fmt.Fprintf(out, "%s%s%s\n", s.dim, r.context(), s.reset)
			}
		}

		if r.cat != prevCat {
			prevCat, prevChk = r.cat, ""
			fmt.Fprintln(out)
			fmt.Fprintf(out, "%s%s%s%s%s  %s%d%s\n", strings.Repeat(" ", bi), s.bold, s.catColor[r.cat], prettyCategoryNames[r.cat], s.reset,
				s.dim, catInObj[r.key+"\x1f"+strconv.Itoa(int(r.cat))], s.reset)
		}

		ind := bi + 2
		pad := strings.Repeat(" ", ind)
		if r.check != prevChk {
			prevChk = r.check
			groupN = chkInObj[r.key+"\x1f"+r.check]
			t := prettyTitleFor(r.check)
			if groupN > 1 {
				t += fmt.Sprintf(" (%d)", groupN)
			}
			s.lr(out, ind, t, s.bold+t+s.reset, r.check, s.cyan+r.check+s.reset)
			shownFix = false
		}

		msg := prettyMessage(r.check, r.msg)
		if groupN > 1 {
			first := "- "
			for _, l := range s.wrap(msg, W-ind-2, 3) {
				fmt.Fprintf(out, "%s%s%s\n", pad, first, l)
				first = "  "
			}
		} else {
			for _, l := range s.wrap(msg, W-ind, 3) {
				fmt.Fprintf(out, "%s%s\n", pad, l)
			}
		}

		if s.showFix && r.rem != "" && !shownFix {
			shownFix = true
			for _, l := range s.wrap("Fix: "+r.rem, W-ind, 2) {
				fmt.Fprintf(out, "%s%s%s%s\n", pad, s.dim, l, s.reset)
			}
		}
	}

	// Footer (kube-linter exits 1 when there are findings)
	fmt.Fprintln(out)
	_, err := fmt.Fprintf(out, "%s%s%s %s%d%s %s in %d %s, exit code 1\n",
		s.red, s.fail, s.reset, s.bold, total, s.reset, fw, nobj, plural(nobj, "object", "objects"))
	return err
}
