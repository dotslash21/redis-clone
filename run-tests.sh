#!/bin/sh
#
# Run all unit and integration tests for the Redis clone
#

set -e # Exit immediately if a command exits with a non-zero status

# Print header
echo "======================================================"
echo "Running all tests for Redis Clone"
echo "======================================================"

# Define colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Function to run tests with nice output
run_tests() {
    local test_path=$1
    local test_name=$2
    
    echo -e "\n${YELLOW}Running $test_name tests...${NC}"
    echo "------------------------------------------------------"
    
    # Run the tests with verbose output and capture result
    if go test -v $test_path; then
        echo -e "\n${GREEN}✓ $test_name tests passed${NC}"
        return 0
    else
        echo -e "\n${RED}✗ $test_name tests failed${NC}"
        return 1
    fi
}

# Track overall success
TESTS_PASSED=true

# Run unit tests in app directory and subdirectories
if ! run_tests "./app/..." "Unit"; then
    TESTS_PASSED=false
fi

echo "\n"

# Run integration tests in tests directory
if ! run_tests "./tests" "Integration"; then
    TESTS_PASSED=false
fi

echo "\n"

# Final summary
if [ "$TESTS_PASSED" = true ]; then
    echo -e "${GREEN}======================================================"
    echo -e "All tests passed successfully!"
    echo -e "======================================================${NC}"
    exit 0
else
    echo -e "${RED}======================================================"
    echo -e "Some tests failed. Please check the output above."
    echo -e "======================================================${NC}"
    exit 1
fi