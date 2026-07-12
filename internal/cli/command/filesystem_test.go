package command
import("testing";"github.com/stretchr/testify/assert";"github.com/Mersad-Moghaddam/syskit/internal/model")
func TestIsPseudo(t *testing.T){assert.True(t,isPseudo("proc"));assert.True(t,isPseudo("pstore"));assert.False(t,isPseudo("ext4"))}
func TestFilesystemTableOptions(t *testing.T){total,free:=uint64(100),uint64(80);table:=filesystemTable(&model.DiskInfo{Mounts:[]model.MountInfo{{MountPoint:"/",FilesystemType:"ext4",Source:"/dev/test",TotalInodes:&total,FreeInodes:&free,Options:[]string{"rw","relatime"}}});assert.Equal(t,"rw,relatime",table.Rows[0][6]);assert.Equal(t,"20%",table.Rows[0][5])}
