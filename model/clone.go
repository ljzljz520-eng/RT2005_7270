package model

func (r Record) Clone() Record   { return r }
func (p Profile) Clone() Profile { return p }
func (a Audit) Clone() Audit {
	m := map[string]string{}
	for k, v := range a.Metadata {
		m[k] = v
	}
	a.Metadata = m
	return a
}
func MergeMetadata(base, extra map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}
