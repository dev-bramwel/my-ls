# Audits Made Easy
Open two terminals. Terminal 1 you will run our binary and terminal 2 will be the system's ls command
## Build binary
```bash
make build
```
### 1. Run both my-ls and the system command ls with no arguments.
```bash
./my-ls
```
### 2. Run both my-ls and the system command ls with the arguments: "file name".
```bash
./my-ls README.md
```
### 3. Run both my-ls and the system command ls with the arguments: "directory name".
```bash
./my-ls docs
```
### 4. Run both my-ls and the system command ls with the flag: "-l".
```bash
./my-ls -l
```
### 5. Run both my-ls and the system command ls with the arguments: "-l file_name".
```bash
./my-ls -l README.md
```
### 6. Run both my-ls and the system command ls with the arguments: "-l <directory name>".
```bash
./my-ls -l docs
```
### 7. Run both my-ls and the system command ls with the flag: "-l /usr/bin".
```bash
./my-ls -l /usr/bin
```
### 8. Run both my-ls and the system command ls with the flag: "-R", in a directory with folders in it.
```bash
./my-ls -R
```
### 9. Run both my-ls and the system command ls with the flag: "-a". 
```bash
./my-ls -a
```
### 10. Run both my-ls and the system command ls with the flag: "-r".
```bash
./my-ls -r
```
### 11. Run both my-ls and the system command ls with the flag: "-t".
```bash
./my-ls -t
```
### 12. Run both my-ls and the system command ls with the flag: "-la".
```bash
./my-ls -la
```
### 13. Run both my-ls and the system command ls with the arguments: "-l -t directory_name".
```bash
./my-ls -l -t docs
```
### 14. Run both my-ls and the system command ls with the arguments: "-lRr directory_name", in which the directory chosen contains folders.
```bash
./my-ls -lRr docs
```
### 15. Run both my-ls and the system command ls with the arguments: "-l directory_name -a file_name"
```bash
./my-ls -l docs -a README.md
```
### 16. Run both my-ls and the system command ls with the arguments: "-lR directory_name///sub_directory_name/// directory_name/sub_directory_name/"
```bash
./my-ls -lR pkg///config/// pkg/config/
```
### 17. Run both my-ls and the system command ls with the arguments: "-la /dev"
```bash
./my-ls -la /dev
```
### 18. Run both my-ls and the system command ls with the arguments: "-alRrt directory_name", in which the directory chosen contains folders and files within folders. Time of modification of all files within that folder must be the same.
```bash
./my-ls -alRrt pkg
```
### 19. Create directory with - name and run both my-ls and the system command ls with the arguments: "-"
```bash
mkdir - && ./my-ls -
```
### 20. Create file and link for this file and run both my-ls-1 and the system command ls with the arguments: "-l symlink_file/"
```bash
# create original file
touch target_file.txt
# create a link to original file
ln -s target_file.txt my_symlink
# Running command
./my-ls -l my_symlink/
```
### 21. Create file and link for this file and run both my-ls-1 and the system command ls with the arguments: "-l symlink_file"
```bash
./my-ls -l my_symlink
```
### 22. Create directory that contains files and link for this directory and run both my-ls-1 and the system command ls with the arguments: "-l symlink_dir"
```bash
# create target dir
mkdir target_dir
# Create files inside the directory
touch target_dir/file1.txt target_dir/file2.txt
# create symlink to the directory
ln -s target_dir my_dir_symlink
# run the system commands
./my-ls -l my_dir_symlink/
```
### 23. Create directory that contains files and link for this directory and run both my-ls-1 and the system command ls with the arguments: "-l <symlink dir>"
```bash
#same as above but without the training backslash
./my-ls -l my_dir_symlink
```
