//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type VerificationCodes struct {
	Vid         int32 `sql:"primary_key"`
	UID         int32
	InstagramID string
	Code        string
}
