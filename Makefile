run:
	source /Users/pascallemerrer/Documents/Dev/projects/bleuvanille/server/env.sh && /Users/pascallemerrer/Documents/Dev/projects/bleuvanille/server/watch.sh


# start docker machine
docker:
	bash --login '/Applications/Docker/Docker Quickstart Terminal.app/Contents/Resources/Scripts/start.sh'


db:
	@echo starts existing ArangoDB container
	docker start arangodb

createdb:
	@echo creates and run a new ArangoDB container
	# allow NFS access from virtual box to host directories
	/usr/local/bin/docker-machine-nfs default
	docker run -e ARANGO_NO_AUTH=1 -d -p 8529:8529 -v /Users/pascallemerrer/Documents/Dev/servers/arangodb/data:/var/lib/arangodb --name=arangodb arangodb/arangodb:2.8.7
	docker ps

 .PHONY: docker db createdb
