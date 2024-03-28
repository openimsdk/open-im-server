// Copyright © 2024 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

const ConfKey = "conf"

const (
	// DefaultDirPerm is used for creating general directories, allowing the owner to read, write, and execute,
	// while the group and others can only read and execute.
	DefaultDirPerm = 0755

	// PrivateFilePerm is used for sensitive files, allowing only the owner to read and write.
	PrivateFilePerm = 0600

	// ExecFilePerm is used for executable files, allowing the owner to read, write, and execute,
	// while the group and others can only read.
	ExecFilePerm = 0754

	// SharedDirPerm is used for shared directories, allowing the owner and group to read, write, and execute,
	// with no permissions for others.
	SharedDirPerm = 0770

	// ReadOnlyDirPerm is used for read-only directories, allowing the owner, group, and others to only read.
	ReadOnlyDirPerm = 0555
)
