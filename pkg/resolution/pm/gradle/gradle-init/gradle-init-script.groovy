def debrickedOutputFile = new File('.debricked.multiprojects.txt')

allprojects {
    task debrickedFindSubProjectPaths() {
        def output = project.name + "\n" + project.projectDir + "\n" + project.rootProject + "\n\n"
        doLast {
            synchronized(debrickedOutputFile) {
                debrickedOutputFile << output
            }
        }
    }
}

allprojects {
    task debrickedAllDeps(type: DependencyReportTask) {
        outputFile = file('./.debricked-gradle-dependencies.txt')
    }
}
