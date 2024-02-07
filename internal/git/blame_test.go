package git

import (
	"io/ioutil"
	"os"
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

	if len(blameRes.Files[0].Lines) == 0 {
		t.Fatal("Should be larger than 0 lines, was", len(blameRes.Files[0].Lines))
	}

	if blameRes.Files[0].Lines[0].Author.Name == "" {
		t.Fatal("Author should not be empty")
	}
	if blameRes.Files[0].Lines[0].Author.Email == "" {
		t.Fatal("Email should not be empty")
	}
}

func TestToFile(t *testing.T) {

	tempFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name()) // clean up

	blamefiles := BlameFiles{
		Files: []BlameFile{
			{
				Lines: []BlameLine{
					{
						Author: Author{
							Email: "example@example.com",
							Name:  "Example",
						},
						LineNumber: 1,
					},
				},
				Path: "example.txt",
			},
		},
	}

	err = blamefiles.ToFile(tempFile.Name())
	if err != nil {
		t.Fatal("failed to write to file. Error:", err)
	}

	content, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := "example.txt,1,Example,example@example.com\n"
	if string(content) != expected {
		t.Errorf("Expected %s, got %s", expected, content)
	}

}
