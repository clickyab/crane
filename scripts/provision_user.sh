#!/bin/bash -x
set -euo pipefail

echo -e "\nexport ENV=development\n" >> /home/develop/.zshrc
echo -e "\nexport PATH=\${PATH}:/home/develop/go/src/clickyab.com/crane/bin" >> /home/develop/.zshrc
echo -e "SAVEHIST=1000" >> /home/develop/.zshrc

#make all
