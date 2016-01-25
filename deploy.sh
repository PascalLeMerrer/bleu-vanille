if [ "$#" -ne 0 ]; then
	echo "deploying Bleu Vanille branch/tag/commit $1" 
	ansible-playbook -i hosts -e reference=$1 playbook.yml
else
	echo "deploying Bleu Vanille branch master" 
	ansible-playbook -i hosts playbook.yml
fi