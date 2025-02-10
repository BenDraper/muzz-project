Setup

make sure you're in the 'muzz-project' directory. Run:

docker-compose build

docker-compose up

This will create two docker containers: a mysql database and an instance of the explore service.

By default this listens on port 8080 but can be changed by updating line 7 in docker-compose.yml to

      - "<desired port>:8080"

e.g.
        - "50051:8080"

ctrl + c will stop both containers.