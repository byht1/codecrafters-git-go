package gitObject

type ObjectCode int

const (
	OBJ_COMMIT_CODE    ObjectCode = 1
	OBJ_TREE_CODE      ObjectCode = 2
	OBJ_BLOB_CODE      ObjectCode = 3
	OBJ_TAG_CODE       ObjectCode = 4
	OBJ_OFS_DELTA_CODE ObjectCode = 6
	OBJ_REF_DELTA_CODE ObjectCode = 7

	TREE_MODE   string = "040000"
	FILE_MODE_1 string = "100644"
	FILE_MODE_2 string = "100755"

	TAG_TYPE    string = "tag"
	BLOB_TYPE   string = "blob"
	TREE_TYPE   string = "tree"
	COMMIT_TYPE string = "commit"

	RAW_TREE_MODE string = "40000"
	PREFIX_MODE   string = "100"
)
