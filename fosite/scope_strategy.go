// Copyright © 2025 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package fosite

import "strings"

// ScopeStrategy is a strategy for matching scopes.
type ScopeStrategy func(haystack []string, needle string) bool

func HierarchicScopeStrategy(haystack []string, needle string) bool {
	for _, this := range haystack {
		// foo == foo -> true
		if this == needle {
			return true
		}

		// picture.read > picture -> false (scope picture includes read, write, ...)
		if len(this) > len(needle) {
			continue
		}

		needles := strings.Split(needle, ".")
		currentHaystack := strings.Split(this, ".")
		haystackLen := len(currentHaystack) - 1
		for k, currentNeedle := range needles {
			if haystackLen < k {
				return true
			}

			current := currentHaystack[k]
			if current != currentNeedle {
				break
			}
		}
	}

	return false
}

func ExactScopeStrategy(haystack []string, needle string) bool {
	for _, this := range haystack {
		if needle == this {
			return true
		}
	}

	return false
}

func WildcardScopeStrategy(matchers []string, needle string) bool {
	needleParts := strings.Split(needle, ".")
	for _, matcher := range matchers {
		matcherParts := strings.Split(matcher, ".")
		if len(matcherParts) > len(needleParts) {
			continue
		}

		var noteq bool
		for k, c := range matcherParts {
			// this is the last item and the lengths are different
			if k == len(matcherParts)-1 && len(matcherParts) != len(needleParts) {
				if c != "*" {
					noteq = true
					break
				}
			}

			if c == "*" && len(needleParts[k]) > 0 {
				// pass because this satisfies the requirements
				continue
			} else if c != needleParts[k] {
				noteq = true
				break
			}
		}

		if !noteq {
			return true
		}
	}

	return false
}

// DeepWildcardScopeStrategy is an opt-in extension of WildcardScopeStrategy.
//
// It preserves existing wildcard behavior and adds:
//   - "*" matches zero or more segments.
//   - "+" matches exactly one non-empty segment.
func DeepWildcardScopeStrategy(matchers []string, needle string) bool {
	for _, pattern := range matchers {
		if deepMatchPattern(pattern, needle) {
			return true
		}
	}

	return false
}

func deepMatchPattern(pattern, candidate string) bool {
	if hasEmptySegment(candidate) {
		return false
	}
	return deepMatchFrom(pattern, 0, candidate, 0)
}

func hasEmptySegment(s string) bool {
	if len(s) == 0 {
		return false
	}
	if s[0] == '.' || s[len(s)-1] == '.' {
		return true
	}
	for i := 1; i < len(s); i++ {
		if s[i] == '.' && s[i-1] == '.' {
			return true
		}
	}
	return false
}

func deepMatchFrom(pattern string, pi int, candidate string, ci int) bool {
	if pi >= len(pattern) && ci >= len(candidate) {
		return true
	}
	if pi >= len(pattern) {
		return false
	}

	pStart := pi
	for pi < len(pattern) && pattern[pi] != '.' {
		pi++
	}
	pEnd := pi
	nextPi := pi
	if nextPi < len(pattern) && pattern[nextPi] == '.' {
		nextPi++
	}

	if pEnd-pStart == 1 && pattern[pStart] == '*' {
		if nextPi >= len(pattern) {
			return true
		}

		ciTry := ci
		for {
			if deepMatchFrom(pattern, nextPi, candidate, ciTry) {
				return true
			}
			if ciTry >= len(candidate) {
				break
			}
			for ciTry < len(candidate) && candidate[ciTry] != '.' {
				ciTry++
			}
			if ciTry < len(candidate) && candidate[ciTry] == '.' {
				ciTry++
			}
		}
		return false
	}

	if ci >= len(candidate) {
		return false
	}
	if candidate[ci] == '.' {
		return false
	}

	cStart := ci
	for ci < len(candidate) && candidate[ci] != '.' {
		ci++
	}
	cEnd := ci
	nextCi := ci
	if nextCi < len(candidate) && candidate[nextCi] == '.' {
		nextCi++
	}

	if pEnd-pStart == 1 && pattern[pStart] == '+' {
		return deepMatchFrom(pattern, nextPi, candidate, nextCi)
	}

	if pEnd-pStart != cEnd-cStart {
		return false
	}
	for k := 0; k < pEnd-pStart; k++ {
		if pattern[pStart+k] != candidate[cStart+k] {
			return false
		}
	}
	return deepMatchFrom(pattern, nextPi, candidate, nextCi)
}
