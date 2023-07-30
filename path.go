package protoparts

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// matches encoded PathTerms and extracts its variables
var (
	//                                 tag ↓      ↓ index?     ↓ key? (hex)
	pathTermRe = regexp.MustCompile(`^(\d+)(?:\[(\d+)])?(?:\[0[xX]([0-9a-fA-F]+)])?$`)
	//                                  field name ↓      ↓ index?     ↓ key? (hex)
	symbolicPathTermRe = regexp.MustCompile(`^(\w+)(?:\[(\d+)])?(?:\[0[xX]([0-9a-fA-F]+)])?$`)
)

// A PathTerm is a single item in a hierarchical Path. It encodes a tag number, and an index denoting which item in a
// repeated sequence the value belongs to.
type PathTerm struct {
	// Tag is the tag number of the field.
	Tag protowire.Number
	// Index is an index into a repeated sequence. -1 indicates the absence of an index.
	Index int
	// Key is the key for the value in a map
	Key []byte
}

// Equal returns whether this term is equivalent to another.
func (t PathTerm) Equal(u PathTerm) bool {
	return t.Tag == u.Tag &&
		t.Index == u.Index &&
		bytes.Equal(t.Key, u.Key)
}

// Compare returns an integer defining the ordering between the two terms. Much as with bytes.Compare, the result is
// 0 if t==u, -1 if t<u, and +1 if t>u.
func (t PathTerm) Compare(u PathTerm) int {
	if t.Tag != u.Tag {
		// t and u are (descendants of) fields of the same message; sort them by their tag numbers
		if t.Tag < u.Tag {
			return -1
		}
		return 1
	}
	if k := bytes.Compare(t.Key, u.Key); k != 0 {
		// t and u are (descendants of) siblings of the same map; sort them lexicographically
		return k
	}
	if t.Index != u.Index {
		// t and u are (descendants of) siblings of the same list; sort them by their indices
		if t.Index < u.Index {
			return -1
		}
		return 1
	}
	return 0 // terms are equal
}

func (t PathTerm) ordinals() string {
	s := ""
	if t.Index != -1 {
		s += fmt.Sprintf("[%d]", t.Index)
	}
	if len(t.Key) > 0 {
		s += fmt.Sprintf("[0x%x]", t.Key)
	}
	return s
}

func (t PathTerm) String() string {
	return fmt.Sprintf("%d%s", t.Tag, t.ordinals())
}

func (t PathTerm) SymbolicString(fd protoreflect.FieldDescriptor) string {
	return fmt.Sprintf("%s%s", fd.Name(), t.ordinals())
}

func decodePathTermOrdinals(idxstr, keystr string) (int, []byte) {
	idx, key := -1, []byte(nil)
	if idxstr != "" {
		_idx, err := strconv.Atoi(idxstr)
		if err != nil {
			return 0, nil
		}
		idx = _idx
	}
	if keystr != "" {
		_key, err := hex.DecodeString(keystr)
		if err != nil {
			return 0, nil
		}
		key = _key
	}
	return idx, key
}

// DecodePathTerm turns an string of the variety returned by PathTerm.String() into a PathTerm. If the term cannot
// be decoded returns a zero PathTerm which can be identified if its Tag == 0 (no valid Protocol Buffer can have a tag
// <1).
//
// If the string is to be constructed by hand, its format is tag[idx][0xkey], where:
//
// - tag is the base-10 representation of the field's tag number
//
// - idx, which is optional, is the base-10 representation of an ordinal within a repeated field, to select an
// individual element from a list
//
// - key, which is optional, is a hex-encoded binary protocol buffer representation prefixed with 0x, of an key within
// a map, to select an individual value from a map
//
// For example: 1 to select the field with number 1, 2[0x0a] to select the key 0a from the map with field number 2,
// or 4[5] to select the fifth item from a list with field number 4.
func DecodePathTerm(s string) PathTerm {
	match := pathTermRe.FindStringSubmatch(s)
	if match == nil {
		return PathTerm{}
	}
	tagstr, idxstr, keystr := match[1], match[2], match[3]
	tag, err := strconv.ParseInt(tagstr, 0, 32)
	if err != nil {
		return PathTerm{} // invalid tag number
	}
	idx, key := decodePathTermOrdinals(idxstr, keystr)
	return PathTerm{
		Tag:   protowire.Number(tag),
		Index: idx,
		Key:   key}
}

// DecodeSymbolicPathTerm turns a symbolically encoded PathTerm as returned by PathTerm.SymbolicString() back into a
// PathTerm. If the term cannot be decoded returns a zero PathTerm which can be identified if its Tag == 0 (no valid
// Protocol Buffer can have a tag <1).
//
// If the string is to be constructed by hand, its format is field[idx][0xkey], where:
//
// - field is the name of the field
//
// - idx, which is optional, is the base-10 representation of an ordinal within a repeated field, to select an
// individual element from a list
//
// - key, which is optional, is a hex-encoded binary protocol buffer representation prefixed with 0x, of an key within
// a map, to select an individual value from a map
//
// For example: foo to select the field with name foo, bar[0x0a] to select the key 0a from the map with field name bar,
// or baz[5] to select the fifth item from a list with field name baz.
func DecodeSymbolicPathTerm(s string, md protoreflect.MessageDescriptor) PathTerm {
	match := symbolicPathTermRe.FindStringSubmatch(s)
	if match == nil {
		return PathTerm{}
	}
	namestr, idxstr, keystr := match[1], match[2], match[3]
	fd := md.Fields().ByName(protoreflect.Name(namestr))
	if fd == nil {
		return PathTerm{}
	}
	idx, key := decodePathTermOrdinals(idxstr, keystr)
	return PathTerm{
		Tag:   fd.Number(),
		Index: idx,
		Key:   key}
}

// A Path is a hierarchical path to a (potentially deeply-nested) value within a serialised Protocol Buffer, identified
// by the tag numbers, and potentially indices and keys, of its ancestors.
//
// A nil Path is different from an empty Path: a nil path addresses nothing, but an empty path addresses the root of a
// message.
type Path []PathTerm

// DecodePath turns an encoded Path of the variety returned by PathTerm.String() into a Path. If the path cannot be
// decoded returns nil.
//
// If the string is to be constructed by hand, its format is a slash-separated concatenation of the strings of its
// terms (see DecodePathTerm).
func DecodePath(s string) Path {
	if s == "" {
		return Path{}
	}
	terms := strings.Split(s, `/`)
	v := make(Path, len(terms))
	for i, term := range terms {
		t := DecodePathTerm(term)
		if t.Tag == 0 {
			return nil // invalid path term
		}
		v[i] = t
	}
	return v
}

// DecodeSymbolicPath turns a symbolically encoded Path (as returned by Path.SymbolicString()) back into a Path. If the
// path cannot be decoded returns nil.
//
// If the string is to be constructed by hand, its format is a slash-separated concatenation of the symbolic strings of\
// its terms (see DecodeSymbolicPathTerm).
func DecodeSymbolicPath(s string, md protoreflect.MessageDescriptor) Path {
	if md == nil {
		return nil
	}
	if s == "" {
		return Path{}
	}
	terms := strings.Split(s, `/`)
	v := make(Path, len(terms))
	for i, term := range terms {
		t := DecodeSymbolicPathTerm(term, md)
		if t.Tag == 0 {
			return nil // invalid path term
		}
		md = md.Fields().ByNumber(t.Tag).Message()
		v[i] = t
	}
	return v
}

// String encodes the Path into a slash-separated string, eg: "1/2/3[4]/5[0x0a]/6[7]". This string can be decoded using
// DecodePath.
//
// The format of this string is suitable for long-term storage and is suitable for equality comparisons.
func (p Path) String() string {
	// WARNING: The output of this string must be kept stable and backwards-compatible: as the comment above explains
	// it is used for storage in the database and can be compared for equality without translation back into a Path.
	// Change it very carefully.
	terms := make([]string, len(p))
	for i, term := range p {
		terms[i] = term.String()
	}
	return strings.Join(terms, `/`)
}

// SymbolicString returns a human-readable representation of the Path with field names in place of tag numbers,
// eg: "foo/bar/baz[4]/foobar[0x0a]/foobaz[7]". This string can be decoded using DecodeSymbolicPath.
//
// Warning: Fields can be renamed freely in Protocol Buffers without affecting their serialisation. So the output of
// this function is NOT guaranteed to be stable, should NOT be stored and is only intended for display to humans. For
// storage, the output of String() should be used since it does not rely on the stability of field names within a
// message descriptor to reconstruct.
func (p Path) SymbolicString(md protoreflect.MessageDescriptor) (string, error) {
	terms := make([]string, len(p))
	for i, term := range p {
		if md == nil {
			return "", errors.New("path not found in message descriptor")
		}
		fd := md.Fields().ByNumber(term.Tag)
		if fd == nil {
			return "", errors.Errorf("path %s not valid in descriptor %s", p, md.FullName())
		}
		terms[i] = term.SymbolicString(fd)
		md = fd.Message()
		// Leave checking whether md is nil until the next iteration both to simplify error handling and because
		// getting a nil descriptor on the last term is not an error
	}
	return strings.Join(terms, `/`), nil
}

// IsValid returns whether this path is valid for the passed message descriptor. This does not necessarily mean that
// lookups are guaranteed to work; this can check the fields used are valid, but if there are list indices it cannot
// know what its cardinality will be.
func (p Path) IsValid(md protoreflect.MessageDescriptor) bool {
	for _, term := range p {
		if md == nil {
			return false
		}
		fd := md.Fields().ByNumber(term.Tag)
		if fd == nil {
			return false
		}
		md = fd.Message()
	}
	return true
}

// Equal returns whether this and the passed Path are equivalent.
func (p Path) Equal(q Path) bool {
	if p == nil || q == nil {
		return p == nil && q == nil
	}
	if len(p) != len(q) {
		return false
	}
	for i := range p {
		if !p[i].Equal(q[i]) {
			return false
		}
	}
	return true
}

// HasPrefix returns whether this Path is a descendant of the prefix. True if the prefix is equal to the Path.
func (p Path) HasPrefix(prefix Path) bool {
	if len(p) < len(prefix) || prefix == nil {
		return false
	}
	for i := range prefix {
		prefixterm, term := prefix[i], p[i]
		if prefixterm.Index == -1 {
			// If the last item of the prefix doesn't specify an index, it doesn't matter whether the last element of the
			// path does. Or in other words: 1.2 is a prefix of 1.2[3], but 1.2[1] is not a prefix of 1.2[3]
			term.Index = -1
		}
		if prefixterm.Key == nil {
			// Similar logic for the key. 1.2 is a prefix of 1.2[0x0a], but 1.2[0x0a] is not a prefix of 1.2[0x0b]
			term.Key = nil
		}
		if !prefixterm.Equal(term) {
			return false
		}
	}
	return true
}

// Parent returns a Path with the last element of this one removed. Eg: 1.2.3 -> 1.2
func (p Path) Parent() Path {
	if len(p) == 0 {
		return p
	}
	return p[:len(p)-1]
}

// Last returns the last element of the path. Eg: 1.2.3 -> 3. If the path is empty, this will panic.
func (p Path) Last() PathTerm {
	return p[len(p)-1]
}

// Compare returns an integer defining the ordering relationship between the two paths. Much as with bytes.Compare,
// the result will be 0 if p==q, -1 if p<q, and +1 if p>q.
func (p Path) Compare(q Path) int {
	// Loop over every element of the path until we find some difference which defines an order between them
	for n := range p {
		if n > len(q)-1 {
			return 1 // p is longer than q but they share a prefix, so p > q
		}
		pterm, qterm := p[n], q[n]
		if r := pterm.Compare(qterm); r != 0 {
			return r
		}
	}
	if len(q) > len(p) {
		return -1 // q is longer than p but they share a prefix, so p < q
	}
	return 0 // paths are equal
}
