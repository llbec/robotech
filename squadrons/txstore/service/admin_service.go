package service

import "txstore/storage/sqlite"

type AdminService struct {
	DBManager *DBManager
}

func NewAdminService(dbm *DBManager) *AdminService {
	return &AdminService{DBManager: dbm}
}

// Split: 手动迁移 active -> 新 archive
func (s *AdminService) SplitMove(fromBlock, toBlock int64, archivePath string) error {
	dst, err := sqlite.Open(archivePath)
	if err != nil {
		return err
	}
	return sqlite.SplitArchiveMove(s.DBManager.Active(), dst, fromBlock, toBlock)
}
