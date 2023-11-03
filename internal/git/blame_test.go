package git

import (
	"testing"
)

func TestBlame(t *testing.T) {

	repo, err := FindRepository("../../.")
	if err != nil {
		t.Fatal("failed to find repo. Error:", err)
	}

	blame := NewBlamer(repo)

	blameRes, err := blame.BlamAllFiles()

	if err != nil {
		t.Fatal("failed to blame file. Error:", err)
	}

	if len(blameRes[0].Lines) == 0 {
		t.Fatal("Should be larger than 0 lines, was", len(blameRes[0].Lines))
	}

	if blameRes[0].Lines[0].Author.Name == "" {
		t.Fatal("Author should not be empty")
	}
	if blameRes[0].Lines[0].Author.Email == "" {
		t.Fatal("Email should not be empty")
	}
}
