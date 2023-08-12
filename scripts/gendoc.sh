#!/usr/bin/env bash
# Copyright © 2023 OpenIM. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Iterates over two directories: 'pkg' and 'internal/pkg'.
for top in pkg internal/pkg
do
    # Finds all subdirectories (including nested ones) under the current directory in the iteration ('pkg' or 'internal/pkg').
    for d in $(find $top -type d)
    do
        # Checks if 'doc.go' doesn't exist in the current subdirectory.
        if [ ! -f $d/doc.go ]; then
            # Checks if there are any '.go' files in the current subdirectory.
            if ls $d/*.go > /dev/null 2>&1; then
                # Echoes the path of the 'doc.go' file to the terminal. 
                # This is likely for debugging or information purposes.
                echo $d/doc.go

                # Writes the package declaration and import comment to the 'doc.go' file in the current subdirectory.
                # 'basename $d' retrieves the name of the current directory (last part of the path).
                # The import comment is constructed based on a static base URL and the directory path.
                echo "package $(basename $d) // import \"github.com/OpenIMSDK/Open-IM-Server/$d\"" > $d/doc.go
            fi
        fi
    done
done
