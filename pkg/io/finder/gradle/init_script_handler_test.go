package gradle

// func TestWriteInitFile(t *testing.T) {
// 	createErr := errors.New("create-error")
// 	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}

// 	sf := InitScriptHandler{fileWriter: fileWriterMock}
// 	err := sf.WriteInitFile()
// 	assert.Equal(t, SetupScriptError{createErr.Error()}, err)

// 	fileWriterMock = &writerTestdata.FileWriterMock{WriteErr: createErr}
// 	sf = InitScriptHandler{initPath: "file", fileWriter: fileWriterMock}
// 	err = sf.WriteInitFile()
// 	assert.Equal(t, SetupScriptError{createErr.Error()}, err)
// }

// func TestWriteInitFileNoInitFile(t *testing.T) {
// 	sf := InitScriptHandler{initPath: "file", fileWriter: nil}
// 	oldGradleInitScript := gradleInitScript
// 	defer func() {
// 		gradleInitScript = oldGradleInitScript
// 	}()
// 	gradleInitScript = embed.FS{}
// 	err := sf.WriteInitFile()
// 	readErr := errors.New("open gradle-init/gradle-init-script.groovy: file does not exist")
// 	assert.Equal(t, SetupScriptError{readErr.Error()}, err)

// }
