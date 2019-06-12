
#!/bin/bash

pushd ..
docker build -t launch-container .

time dgoss run -v "${HOME}/.aws/credentials:/app/credentials" launch-container || exit 1
popd

exit 0
