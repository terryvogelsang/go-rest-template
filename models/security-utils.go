package models

// Custom types to make sure users cannot edit sensitive properties through API Calls

// ReadOnlyString : String that cannot be set on JSON Unmarshalling
type ReadOnlyString string

// UnmarshalJSON : Make sure ReadOnlyString cannot be set on JSON Unmarshalling
func (ReadOnlyString) UnmarshalJSON([]byte) error { return nil }
