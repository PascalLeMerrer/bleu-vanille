
if [ "$#" -ne 0 ]; then
	echo "deploying Bleu Vanille branch/tag/commit $1" 
	ref=$1
else
	echo "deploying Bleu Vanille branch master" 
	ref="master"
fi
ansible-playbook -i hosts -e reference=$ref --tags "checkout" playbook.yml

cd "./deploy/"
./node_modules/.bin/cucumber.js

ansible-playbook -i hosts --tags "package,supervisor,deploy" playbook.yml