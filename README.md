# Explore Service

# Setup

make sure you're in the 'muzz-project' directory. Run:

`docker-compose build`

`docker-compose up`

This will create two docker containers: a mysql database and an instance of the explore service.

By default this grpc service listens on port 8080 but can be changed by updating line 7 in docker-compose.yml to
        `- "desired port:8080"`

e.g.
        `- "50051:8080"`

ctrl + c will stop both containers.

## Testing

Running

`go test ./...`

Will run all tests. 

You may need to revendor if dependencies haven't been downloaded. This can be done by running

`go mod tidy`
`go mod vendor`

main_test.go is an integration test which spins up a containerised db for the tests.

## Notes

I chose pagination rather than streaming for listing likes because I think it fits the use case better. 
I assumed the most common use case will be for the user to see likes in their feed so having all of these 
on the client device at once would be unnecessary. The choice of 1000 as the maximum responses was fairly arbitrary. 
In practice I'd discuss this decision with others.

I decided to create a mock of the inerface created bythe proto so that downstream users can write nicer unit tests for 
their services that call this service. They can use this pre-generated mock rather than having to stub it themselves.

I used mysql because this data looks fairly structured and because you have mysql listed in your tech stack.