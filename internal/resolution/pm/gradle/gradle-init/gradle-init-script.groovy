def debrickedOutputFile = new File('.debricked.multiprojects.txt')

allprojects {
    task debrickedFindSubProjectPaths() {
        String output = project.projectDir 
        doLast {
            synchronized(debrickedOutputFile) {
                debrickedOutputFile << output + System.getProperty("line.separator")
            }
        }
    }
}

allprojects {
    task debrickedAllDeps(type: DependencyReportTask) {
        outputFile = file('./gradle.debricked.lock')
    }
}
