docker image build -f Dockerfile -t forum .
docker container run -p 8080:8080 --detach --name form forum  