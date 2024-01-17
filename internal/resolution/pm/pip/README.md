# Pip resolution logic

The way resolution of pip lock files works is as follows:

1. Create a Venv in which we do the installation and run all commands
2. Run `pip install -r <requirements.txt_file>` in order to install all dependencies
3. Run `cat` to get the contents of the requirements.txt file
4. Run `pip list` to get a list of all installed packages
5. Run `pip show <list_of_installed_packages>` to get more in-depth information from each package, including the relations between dependencies

The results of the commands above are then combined to form the finished lock file with the following sections:

1. The contents of the requirements.txt (from cat)
2. The list of all installed dependencies (from pip list)
3. More detailed information on each package with relations (from pip show)
