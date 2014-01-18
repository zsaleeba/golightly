package golightly

import (
	"testing"
	"errors"
	"fmt"
)

//test function starts with "Test" and takes a pointer to type testing.T
func TestTokenList(t *testing.T) {
	tl := NewTokenList("-")

	// add some tokens
	tl.Add(SrcLoc{1, 0}, TokenPackage)
	tl.AddString(SrcLoc{1, 8}, TokenIdentifier, "golightly")
	tl.Add(SrcLoc{3, 0}, TokenImport)
	tl.Add(SrcLoc{3, 7}, TokenOpenGroup)
	tl.AddString(SrcLoc{4, 4}, TokenIdentifier, "testing")
	tl.Add(SrcLoc{5, 0}, TokenCloseGroup)
	tl.AddString(SrcLoc{7, 0}, TokenIdentifier, "i")
	tl.Add(SrcLoc{7, 2}, TokenDeclareAssign)
	tl.AddUInt(SrcLoc{7, 4}, TokenUint, 42)
	tl.Add(SrcLoc{7, 6}, TokenSemicolon)
	tl.AddString(SrcLoc{8, 0}, TokenIdentifier, "j")
	tl.Add(SrcLoc{8, 2}, TokenDeclareAssign)
	tl.AddFloat(SrcLoc{8, 4}, 7.2)
	tl.Add(SrcLoc{8, 7}, TokenSemicolon)

	// now try to get them back out
	tl.StartReading()

	err := checkToken(tl, 1, 0, TokenPackage)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenString(tl, 1, 8, TokenIdentifier, "golightly")
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 3, 0, TokenImport)
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 3, 7, TokenOpenGroup)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenString(tl, 4, 4, TokenIdentifier, "testing")
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 5, 0, TokenCloseGroup)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenString(tl, 7, 0, TokenIdentifier, "i")
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 7, 2, TokenDeclareAssign)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenUint(tl, 7, 4, TokenUint, 42)
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 7, 6, TokenSemicolon)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenString(tl, 8, 0, TokenIdentifier, "i")
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 8, 2, TokenDeclareAssign)
	if err != nil {
		t.Error(err)
	}

	err = checkTokenFloat(tl, 8, 4, TokenFloat32, 7.2)
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 8, 6, TokenSemicolon)
	if err != nil {
		t.Error(err)
	}

	err = checkToken(tl, 8, 7, TokenEndOfSource)
	if err != nil {
		t.Error(err)
	}
}

func checkToken(tl *TokenList, line int, column int, tok Token) error {
	var loc SrcLoc
	foundToken := tl.GetToken(&loc)
	if foundToken != tok {
		return errors.New(fmt.Sprint("wrong token: ", foundToken, " != ", tok))
	}

	if loc.Line != line {
		return errors.New(fmt.Sprint("wrong line: ", loc.Line, " != ", line))
	}

	if loc.Column != column {
		return errors.New(fmt.Sprint("wrong column: ", loc.Column, " != ", column))
	}

	return nil
}

func checkTokenString(tl *TokenList, line int, column int, tok Token, str string) error {
	err := checkToken(tl, line, column, tok)
	if err != nil {
		return err
	}

	if tl.GetValueString() != str {
		return errors.New(fmt.Sprint("wrong string: '", tl.GetValueString(), "' != '", str, "'"))
	}

	return nil
}

func checkTokenUint(tl *TokenList, line int, column int, tok Token, v uint64) error {
	err := checkToken(tl, line, column, tok)
	if err != nil {
		return err
	}

	if tl.GetValueUint64() != v {
		return errors.New(fmt.Sprint("wrong uint: ", tl.GetValueUint64(), " != ", v))
	}

	return nil
}

func checkTokenFloat(tl *TokenList, line int, column int, tok Token, v float64) error {
	err := checkToken(tl, line, column, tok)
	if err != nil {
		return err
	}

	if tl.GetValueFloat64() != v {
		return errors.New(fmt.Sprint("wrong float: ", tl.GetValueFloat64(), " != ", v))
	}

	return nil
}
