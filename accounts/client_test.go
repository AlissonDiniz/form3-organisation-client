package accounts

import (
	"fmt"
	"testing"
	"strings"
	"github.com/bxcodec/faker/v3"
)

var Countries = [5]string{"GB", "AU", "CA", "DE", "PT"}

type AccountDataSeed struct {
	Country  int   `faker:"boundary_start=0, boundary_end=4"`
	BankIDCode  string   `faker:"len=16"`
	Bic1  string   `faker:"len=6"`
	Bic2  int   `faker:"boundary_start=10, boundary_end=99"`
	SecondaryIdentification  string   `faker:"len=10"`
}

func createAccount() *AccountData {
	account_seed := AccountDataSeed{}
	_ = faker.FakeData(&account_seed)

	account_classification := "Personal"
	joint_account := false
	account_matching_opt_out := false
	attributes := AccountAttributes {
		Country: &Countries[account_seed.Country],
		BaseCurrency: faker.Currency(),
		BankID: faker.CCNumber(),
		BankIDCode: strings.ToUpper(account_seed.BankIDCode),
		Bic: fmt.Sprintf("%s%d", strings.ToUpper(account_seed.Bic1), account_seed.Bic2),
		Name: []string{faker.Name()},
		AlternativeNames: []string{faker.Name()},
		AccountClassification: &account_classification,
		JointAccount: &joint_account,
		AccountMatchingOptOut: &account_matching_opt_out,
		SecondaryIdentification: strings.ToUpper(account_seed.SecondaryIdentification),
	}

	account_data := AccountData{
		Type: "accounts",
		ID: faker.UUIDHyphenated(),
		OrganisationID: faker.UUIDHyphenated(),
		Attributes: &attributes,
	}
	return &account_data
}

// TestCreateAccount calls form3.com/organization-api/accounts.Create with valid values, checking
func TestCreateAccount(t *testing.T) {
	account_data := createAccount()
	
	err := Create(account_data)
	if err != nil {
		t.Fatalf("This test should not return an error");
	}
}

// TestCreateAccountWithConflict calls form3.com/organization-api/accounts.Create with duplicated account, checking
func TestCreateAccountWithConflict(t *testing.T) {
	account_data := createAccount()
	Create(account_data)
	
	err := Create(account_data)
	if err != nil {
		switch err.(type) {
			case *ConflictError :
			default:
        t.Fatalf("This test must return an error ConflictError");
    }
	} else {
		t.Fatalf("This test must return an error ConflictError");
	}
}

// TestFetchAccountNotFound calls form3.com/organization-api/accounts.Fetch with an Id that does not exist, checking
func TestFetchAccountNotFound(t *testing.T) {
	id := faker.UUIDHyphenated()

	_, err := Fetch(id)
	if err != nil {
		switch err.(type) {
			case *NotFoundError :
			default:
        t.Fatalf("This test must return an error NotFoundError");
    }
	} else {
		t.Fatalf("This test must return an error NotFoundError");
	}
}

// TestFetchAccount calls form3.com/organization-api/accounts.Fetch with an existing id, checking
func TestFetchAccount(t *testing.T) {
	account_data := createAccount()
	Create(account_data)

	res, err := Fetch(account_data.ID)
	if err != nil {
		t.Fatalf("This test should not return an error");
	}
	
	if account_data.ID != res.ID {
		t.Fatalf("Account ID is different to Account ID fetched");
	} 
}

// TestDeleteAccountNotFoundInvalidId calls form3.com/organization-api/accounts.Delete with an Id that does not exist, checking
func TestDeleteAccountNotFoundInvalidId(t *testing.T) {
	id := faker.UUIDHyphenated()
	version := 0

	err := Delete(id, version)
	if err != nil {
		switch err.(type) {
			case *NotFoundError :
			default:
        t.Fatalf("This test must return an error NotFoundError");
    }
	} else {
		t.Fatalf("This test must return an error NotFoundError");
	}
}

// TestDeleteAccountNotFoundInvalidVersion calls form3.com/organization-api/accounts.Delete with a version that does not exist, checking
func TestDeleteAccountNotFoundInvalidVersion(t *testing.T) {
	id := faker.UUIDHyphenated()
	version := 10

	err := Delete(id, version)
	if err != nil {
		switch err.(type) {
			case *NotFoundError :
			default:
        t.Fatalf("This test must return an error NotFoundError");
    }
	} else {
		t.Fatalf("This test must return an error NotFoundError");
	}
}

// TestDeleteAccount calls form3.com/organization-api/accounts.Delete with an existing id, checking
func TestDeleteAccount(t *testing.T) {
	account_data := createAccount()
	Create(account_data)
	version := 0

	err := Delete(account_data.ID, version)
	if err != nil {
		t.Fatalf("This test should not return an error");
	}
}