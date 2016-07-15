#!/bin/bash
SECONDS=0

if [ "$#" -ne 0 ]; then
	echo "deploying Bleu Vanille branch/tag/commit $1"
	ref=$1
else
	echo "deploying Bleu Vanille branch master"
	ref="master"
fi
ansible-playbook -i hosts -e reference=$ref --ask-become-pass --tags "checkout" playbook.yml

cd ./deploy/

source ../env.sh && ./server &

if [ "$?" -ne 0 ]; then
    echo "FATAL. Build failed. Deployment canceled."
    exit $?
fi

processId=$!

../node_modules/.bin/cucumber.js --fail-fast

if [ "$?" -ne 0 ]; then
	echo "FATAL. Test execution failed. Deployment canceled."
	kill -9 $processId
	exit $?
fi

kill -9 $processId



cd ..
ansible-playbook -i hosts --ask-become-pass -vvvv --tags "package,supervisor,deploy" playbook.yml

echo "Release deployed in $SECONDS seconds"