package chunk

type ChunkServer struct {
	RootDir string // Root directory to store blocks
}

// NewChunkServer returns a server which store data.
func NewChunkServer(rootDir string) *ChunkServer {
	c := &ChunkServer{
		RootDir:  rootDir,
	}
	return c
}