package entity

import "strconv"

type Dentist struct {
	Name string
	CRO  string
}

// ValidateCRO checks if the CRO field is a valid integer
func (d *Dentist) ValidateCRO() bool {
	_, err := strconv.Atoi(d.CRO)
	return err == nil
}

// ValidateName checks if the Name field has a maximum length of 10 characters
func (d *Dentist) ValidateName() bool {
	return len(d.Name) <= 10
}
