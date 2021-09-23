if [[ ! -d "e2etests/test_helper/bats-support" ]]; then
  # Download bats-support dynamically so it doesnt need to be added into source
  git clone https://github.com/ztombol/bats-support e2etests/test_helper/bats-support --depth 1
fi

if [[ ! -d "e2etests/test_helper/redhatcop-bats-library" ]]; then
  # Download redhat-cop/bats-library dynamically so it doesnt need to be added into source
  git clone https://github.com/redhat-cop/bats-library e2etests/test_helper/redhatcop-bats-library --depth 1
fi