#!/bin/bash

# Delete temporary transfer files after copying
rm -f data_complement_model_for_common.go
rm -f data_complement_model_common.go
rm -f create_models_dir.sh

echo "Temporary files cleaned up successfully."
echo ""
echo "Note: The original file at features/data_complement/models/data_complement_model.go"
echo "      can be deleted once all code is updated to use the common version."
