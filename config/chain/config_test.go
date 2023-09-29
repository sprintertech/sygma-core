package chain

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	var id uint8 = 1
	valid := ChainConfig{
		Name:     "chain",
		Id:       &id,
		Endpoint: "endpoint",
	}

	missingEndpoint := ChainConfig{
		Name:     "chain",
		Id:       &id,
		Endpoint: "",
	}

	missingName := ChainConfig{
		Name:     "",
		Id:       &id,
		Endpoint: "endpoint",
	}

	missingId := ChainConfig{
		Name:     "chain",
		Endpoint: "endpoint",
	}

	err := valid.Validate()
	if err != nil {
		t.Fatal(err)
	}

	err = missingEndpoint.Validate()
	if err == nil {
		t.Fatalf("must require endpoint field, %v", err)
	}

	err = missingName.Validate()
	if err == nil {
		t.Fatal("must require name field")
	}

	err = missingId.Validate()
	if err == nil {
		t.Fatalf("must require domain id field, %v", err)
	}
}
