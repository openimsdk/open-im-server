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
