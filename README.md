# Clickyab exchange project 

## How to setup this for development?

make sure the folder is cloned in proper folder in your system. this is 
not required for building the project, but for your IDE, its required. 

the project folder must named `exchange` and must be in a folder named `src/clickyab.com`
so the folder is typically like this : 

`
~/project/src/clickyab.com/exchange
`

the go-path in the above system must be `GOPATH=~/project/` but for building it, 
you must run vagrant : 

```
vagrant up
vagrant ssh 
cdp #cdp is an alias for cd into the prject inside the vagrant box
make
```

## How to add new dependency 

The entire process is in your system, or else by destroy the box, you need to download the 
entire dependencies each time
 - install the latest version of `godep` from this [repository](https://github.com/tools/godep)
 - cd into project folder and make sure `GOPATH` is set to the parent folder (see above)
 - use go get (or simply clone the repo in the correct folder in go-path)
 - then run `godep save ./...`
 - Commit the change in one single commit. no other change is allowed in that commit
 
 
## How to contribute

 - Each commit must done one task, or less. its not allowed to do TWO task in one commit, 
 its allowed to done a task into multiple commit
 - Commit message mus follow this role:
    ```
     commit message, desriptive, not longer than 80 char (required)
     
     commit detail after a single empty line. detail are going here, any line 
     is allowed. (optional)
     follow this must be the task refs or fix commit (optional)
     fix #9999 
     refs #222
    ```
 - If a task is done in multiple part, one create a fork, and start the task
 others use PR to that fork to add commits to that fork then PR the result into the 
 main repository. 
 Its very important to use proper message for your commit. describe content of the 
 commit, and make sure the message is not about peoples, or any unrelated things

