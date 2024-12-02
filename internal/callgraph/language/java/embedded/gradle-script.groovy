

allprojects{
    task debrickedCopyDependencies(type: Copy) {
        into ".debrickedTmpDir"
        from configurations.default
    }
}