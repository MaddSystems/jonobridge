#!/bin/bash

# Navigate to the tests directory
cd "$(dirname "$0")"

# Path to the virtual environment
VENV_PATH="../frontend/.venv/bin/activate"

if [ ! -f "$VENV_PATH" ]; then
    echo "Error: Virtual environment not found at $VENV_PATH"
    exit 1
fi

echo "Activating Python Virtual Environment..."
source "$VENV_PATH"

echo "Running Integration Audit Test..."
OUTPUT_FILE="integration_test_result.log"

# Run the python test, capture output to file and display on screen
python3 integration_audit_v2.py 2>&1 | tee "$OUTPUT_FILE"

echo "---------------------------------------------------"
echo "Test execution finished."
echo "Full results saved to: $(pwd)/$OUTPUT_FILE"
