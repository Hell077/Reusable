package crypto

import "testing"

//90% coverage

func TestEncrypt(t *testing.T) {
	testsPass := []string{
		"7Gk9&dKs3",
		"Aq2^Mnz@5",
		"pL9*Dsd!3",
		"Zx8!vBq1#",
		"bC5#Kr$8d",
		"",
	}

	for _, v := range testsPass {
		str, err := Encrypt(v)
		if v == "" && err == nil {
			t.Error("expected error for empty password, got nil")
		}
		if err != nil && v != "" {
			t.Errorf("unexpected error for password %s: %v", v, err)
		}
		if v != "" && str == "" {
			t.Error("encrypt fail, hash is empty")
		}
		if str == v && v != "" {
			t.Errorf("error encrypting password, password equals hash: pass=%s hash=%s", v, str)
		}
	}
}

func TestDecrypt(t *testing.T) {
	testsHash := []string{
		"$2a$10$QmFIZH.i9U4c7YEEGdc8L.IXVRGj1rqZyjFn/9Q4GLSbSxiY/opTm",
		"$2a$10$kIomDC7iKvGX1.yvKu9xG.EBRdrOzvsF5vQF/SOvjTC3UtfAbz.b6",
		"$2a$10$kTSvm9xiDmWLsk9DtgYtx.67cvgLQZQ9LFlRqIKDsI6EEj1aYgBGa",
		"$2a$10$iC9JNSLQOPMnzX76fP82fOtIGdnnZnYoUPzibuODQ4dVym/FiNrAO",
		"$2a$10$2vvmsdef7/nOB.UoaBjdV.U87s6W8XR.pxHlQaKpTA.uXJNKCNwlC",
		"",
	}
	testsPass := []string{
		"7Gk9&dKs3",
		"Aq2^Mnz@5",
		"pL9*Dsd!3",
		"Zx8!vBq1#",
		"bC5#Kr$8d",
		"wrongpass",
	}

	for i := 0; i < len(testsHash); i++ {
		res, err := Decrypt(testsHash[i], testsPass[i])
		if testsPass[i] == "" && err == nil {
			t.Error("expected error for empty password, got nil")
		}
		if testsHash[i] == "" && err == nil {
			t.Error("expected error for empty hash, got nil")
		}
		if err != nil && testsPass[i] != "" && testsHash[i] != "" {
			t.Errorf("unexpected error for hash %s and password %s: %v", testsHash[i], testsPass[i], err)
		}
		if res == false && testsPass[i] != "wrongpass" {
			t.Errorf("decrypt failed for valid password %s", testsPass[i])
		}
		if res == true && testsPass[i] == "wrongpass" {
			t.Error("decrypt succeeded for wrong password")
		}
	}
}
