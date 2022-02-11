package gorm

import (
	"strings"
)

// QuerySignature returns the "signature" for a query:
// a high level description of the operation.
//
// For DDL statements (CREATE, DROP, ALTER, etc.), we we only
// report the first keyword, on the grounds that these statements
// are not expected to be common within the hot code paths of
// an application. For SELECT, INSERT, and UPDATE, and DELETE,
// we attempt to extract the first table name. If we are unable
// to identify the table name, we simply omit it.
func QuerySignature(query string) string {
	s :=NewScanner(query)
	for s.Scan() {
		if s.Token() != COMMENT {
			break
		}
	}

	scanUntil := func(until Token) bool {
		for s.Scan() {
			if s.Token() == until {
				return true
			}
		}
		return false
	}
	scanToken := func(tok Token) bool {
		for s.Scan() {
			switch s.Token() {
			case tok:
				return true
			case COMMENT:
			default:
				return false
			}
		}
		return false
	}

	switch s.Token() {
	case CALL:
		if !scanUntil(IDENT) {
			break
		}
		return "CALL " + s.Text()

	case DELETE:
		if !scanUntil(FROM) {
			break
		}
		if !scanToken(IDENT) {
			break
		}
		tableName := s.Text()
		for scanToken(PERIOD) && scanToken(IDENT) {
			tableName += "." + s.Text()
		}
		return "DELETE FROM " + tableName

	case INSERT, REPLACE:
		action := s.Text()
		if !scanUntil(INTO) {
			break
		}
		if !scanToken(IDENT) {
			break
		}
		tableName := s.Text()
		for scanToken(PERIOD) && scanToken(IDENT) {
			tableName += "." + s.Text()
		}
		return action + " INTO " + tableName

	case SELECT:
		var level int
	scanLoop:
		for s.Scan() {
			switch tok := s.Token(); tok {
			case LPAREN:
				level++
			case RPAREN:
				level--
			case FROM:
				if level != 0 {
					continue scanLoop
				}
				if !scanToken(IDENT) {
					break scanLoop
				}
				tableName := s.Text()
				for scanToken(PERIOD) && scanToken(IDENT) {
					tableName += "." + s.Text()
				}
				return "SELECT FROM " + tableName
			}
		}

	case UPDATE:
		// Scan for the table name. Some dialects allow
		// option keywords before the table name.
		var havePeriod, haveFirstPeriod bool
		if !scanToken(IDENT) {
			return "UPDATE"
		}
		tableName := s.Text()
		for s.Scan() {
			switch tok := s.Token(); tok {
			case IDENT:
				if havePeriod {
					tableName += s.Text()
					havePeriod = false
				}
				if !haveFirstPeriod {
					tableName = s.Text()
				} else {
					// Two adjacent identifiers found
					// after the first period. Ignore
					// the secondary ones, in case they
					// are unknown keywords.
				}
			case PERIOD:
				haveFirstPeriod = true
				havePeriod = true
				tableName += "."
			default:
				return "UPDATE " + tableName
			}
		}
	}

	// If all else fails, just return the first token of the query.
	fields := strings.Fields(query)
	if len(fields) == 0 {
		return ""
	}
	return strings.ToUpper(fields[0])
}
