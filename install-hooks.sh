#!/bin/bash

for hook in $(ls -1 git/hooks/*); do
    if [ -z $hook ]; then
        echo "#!/bin/bash" > ${$(basename git/hooks/pre-commit.old)%.*}
    fi
    name=$(basename $hook)
    echo "bash ./$hook" >> .git/hooks/${name%.*}
done