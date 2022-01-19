package runescape

import "testing"

func TestFormatB(t *testing.T) {
	price := RSPrice("2.0b")
	err := price.Format()
	if err != nil {
		t.Errorf("format failed: %v", err)
	}
	if price != RSPrice("2.0b") {
		t.Errorf("invalid price change, new price: %v", price)
	}
}

func TestFormatM(t *testing.T) {
	price := RSPrice("25.0m")
	err := price.Format()
	if err != nil {
		t.Errorf("format failed: %v", err)
	}
	if price != RSPrice("25.0m") {
		t.Errorf("invalid price change, should be 25.0m but is %v", price)
	}

	price = RSPrice("1252.5m")
	err = price.Format()
	if err != nil {
		t.Errorf("format failed: %v", err)
	}
	if price != RSPrice("1.2b") {
		t.Errorf("invalid price change, should be 1.2b but is %v", price)
	}
}

func TestFormatK(t *testing.T) {
	price := RSPrice("125.3k")
	err := price.Format()
	if err != nil {
		t.Errorf("format failed: %v", err)
	}
	if price != RSPrice("125.3k") {
		t.Errorf("invalid price change, should be 125.3k but is %v", price)
	}

	price = RSPrice("1252.5k")
	err = price.Format()
	if err != nil {
		t.Errorf("format failed: %v", err)
	}
	if price != RSPrice("1.2m") {
		t.Errorf("invalid price change, should be 1.2m but is %v", price)
	}

	price = RSPrice("1252236.5k")
	err = price.Format()
	if err != nil {
		t.Errorf("format failed: %v", err)
	}
	if price != RSPrice("1.2b") {
		t.Errorf("invalid price change, should be 1.2b but is %v", price)
	}
}

func TestFormat(t *testing.T) {
	price := RSPrice("1235")
	err := price.Format()
	if err != nil {
		t.Errorf("format failed: %v", err)
	}
	if price != RSPrice("1235") {
		t.Errorf("invalid price change, should be 1235 but is %v", price)
	}

	price = RSPrice("1235.0000")
	err = price.Format()
	if err != nil {
		t.Errorf("format failed: %v", err)
	}
	if price != RSPrice("1235") {
		t.Errorf("invalid price change, should be 1235 but is %v", price)
	}

	price = RSPrice("12500")
	err = price.Format()
	if err != nil {
		t.Errorf("format failed: %v", err)
	}
	if price != RSPrice("12.5k") {
		t.Errorf("invalid price change, should be 12.5k but is %v", price)
	}

}

func TestCompareGt(t *testing.T) {
	bigger, smaller := RSPrice("5.3b"), RSPrice("123456")
	got, err := bigger.Compare(smaller)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 1 {
		t.Errorf("invalid comparison, should be 1 but is %v", got)
	}

	bigger, smaller = RSPrice("5.3b"), RSPrice("5b")
	got, err = bigger.Compare(smaller)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 1 {
		t.Errorf("invalid comparison, should be 1 but is %v", got)
	}

	bigger, smaller = RSPrice("5.3m"), RSPrice("54562")
	got, err = bigger.Compare(smaller)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 1 {
		t.Errorf("invalid comparison, should be 1 but is %v", got)
	}

	bigger, smaller = RSPrice("5.3m"), RSPrice("5m")
	got, err = bigger.Compare(smaller)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 1 {
		t.Errorf("invalid comparison, should be 1 but is %v", got)
	}

	bigger, smaller = RSPrice("5.3k"), RSPrice("1256")
	got, err = bigger.Compare(smaller)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 1 {
		t.Errorf("invalid comparison, should be 1 but is %v", got)
	}

	bigger, smaller = RSPrice("5.3k"), RSPrice("5.2k")
	got, err = bigger.Compare(smaller)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 1 {
		t.Errorf("invalid comparison, should be 1 but is %v", got)
	}

	bigger, smaller = RSPrice("1255"), RSPrice("256")
	got, err = bigger.Compare(smaller)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 1 {
		t.Errorf("invalid comparison, should be 1 but is %v", got)
	}
}

func TestCompareLt(t *testing.T) {
	bigger, smaller := RSPrice("5.3b"), RSPrice("123456")
	got, err := smaller.Compare(bigger)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != -1 {
		t.Errorf("invalid comparison, should be -1 but is %v", got)
	}

	bigger, smaller = RSPrice("5.3b"), RSPrice("5b")
	got, err = smaller.Compare(bigger)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != -1 {
		t.Errorf("invalid comparison, should be -1 but is %v", got)
	}

	bigger, smaller = RSPrice("5.3m"), RSPrice("54562")
	got, err = smaller.Compare(bigger)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != -1 {
		t.Errorf("invalid comparison, should be -1 but is %v", got)
	}

	bigger, smaller = RSPrice("5.3m"), RSPrice("5m")
	got, err = smaller.Compare(bigger)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != -1 {
		t.Errorf("invalid comparison, should be -1 but is %v", got)
	}

	bigger, smaller = RSPrice("5.3k"), RSPrice("1256")
	got, err = smaller.Compare(bigger)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != -1 {
		t.Errorf("invalid comparison, should be -1 but is %v", got)
	}

	bigger, smaller = RSPrice("5.3k"), RSPrice("5.2k")
	got, err = smaller.Compare(bigger)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != -1 {
		t.Errorf("invalid comparison, should be -1 but is %v", got)
	}

	bigger, smaller = RSPrice("1255"), RSPrice("256")
	got, err = smaller.Compare(bigger)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != -1 {
		t.Errorf("invalid comparison, should be -1 but is %v", got)
	}
}

func TestCompareEq(t *testing.T) {
	p1, p2 := RSPrice("2.0b"), RSPrice("2.0b")
	got, err := p1.Compare(p2)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 0 {
		t.Errorf("invalid comparison, should be 0 got is %v", got)
	}

	p1, p2 = RSPrice("2.0b"), RSPrice("2000m")
	got, err = p1.Compare(p2)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 0 {
		t.Errorf("invalid comparison, should be 0 got is %v", got)
	}

	p1, p2 = RSPrice("2.0m"), RSPrice("2m")
	got, err = p1.Compare(p2)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 0 {
		t.Errorf("invalid comparison, should be 0 got is %v", got)
	}

	p1, p2 = RSPrice("2.0m"), RSPrice("2000k")
	got, err = p1.Compare(p2)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 0 {
		t.Errorf("invalid comparison, should be 0 got is %v", got)
	}

	p1, p2 = RSPrice("2.0k"), RSPrice("2k")
	got, err = p1.Compare(p2)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 0 {
		t.Errorf("invalid comparison, should be 0 got is %v", got)
	}

	p1, p2 = RSPrice("2.0k"), RSPrice("2000")
	got, err = p1.Compare(p2)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 0 {
		t.Errorf("invalid comparison, should be 0 got is %v", got)
	}

	p1, p2 = RSPrice("1235"), RSPrice("1235")
	got, err = p1.Compare(p2)
	if err != nil {
		t.Errorf("comparison failed: %v", err)
	}
	if got != 0 {
		t.Errorf("invalid comparison, should be 0 got is %v", got)
	}
}

func TestMultiply(t *testing.T) {
	price, multiplier := RSPrice("125"), 100
	err := price.Multiply(multiplier)
	if err != nil {
		t.Errorf("multplication failed: %v", err)
	}
	if price != RSPrice("12.5k") {
		t.Errorf("invalid multiplication, should be 12.5k but got %v", price)
	}

	price, multiplier = RSPrice("125.2m"), 10
	err = price.Multiply(multiplier)
	if err != nil {
		t.Errorf("multplication failed: %v", err)
	}
	if price != RSPrice("1.2b") {
		t.Errorf("invalid multiplication, should be 1.2b but got %v", price)
	}
}
